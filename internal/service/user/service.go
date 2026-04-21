package user

import (
	authdto "back/internal/dto/auth"
	userdto "back/internal/dto/user"
	mail "back/internal/mailer"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"back/internal/cache"
	"back/internal/model"
	"back/internal/repository"
	jwtpkg "back/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginAlreadyExists      = errors.New("login already exists")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrGoogleOnlyAuth          = errors.New("google only auth")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrEmailDeliveryFailed     = errors.New("email delivery failed")
)

type Service struct {
	userRepo   repository.UserRepository
	jwtService *jwtpkg.Service
	store      *cache.MemoryStore
	mailer     *mail.SMTPMailer
}

func NewUserService(
	userRepo repository.UserRepository,
	jwtSecret string,
	store *cache.MemoryStore,
	mailer *mail.SMTPMailer,
) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtService: jwtpkg.NewService(jwtSecret),
		store:      store,
		mailer:     mailer,
	}
}

func (s *Service) Register(login, email, password string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	login = strings.TrimSpace(login)

	existing, err := s.userRepo.GetByLoginOrEmail(login, email)
	if err != nil {
		return err
	}

	if existing != nil {
		if existing.Login == login {
			return ErrLoginAlreadyExists
		}
		if existing.Email == email {
			return ErrEmailAlreadyExists
		}
	}

	code := generateCode()
	pending := authdto.PendingRegistration{
		Login:     login,
		Email:     email,
		Password:  password,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	payload, err := json.Marshal(pending)
	if err != nil {
		return err
	}

	subject := "Код подтверждения регистрации"
	if err := s.mailer.Send(email, subject, code); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailDeliveryFailed, err)
	}

	if err := s.store.Set(pendingRegistrationKey(email), string(payload)); err != nil {
		return err
	}

	return nil
}

func (s *Service) CompleteVerification(body userdto.VerifyRequest) (string, error) {
	email := strings.TrimSpace(strings.ToLower(body.Email))
	code := strings.TrimSpace(body.Code)

	raw, err := s.store.Get(pendingRegistrationKey(email))
	if err != nil {
		return "", ErrInvalidVerificationCode
	}

	var pending authdto.PendingRegistration
	if err := json.Unmarshal([]byte(raw), &pending); err != nil {
		return "", err
	}

	if pending.Code != code {
		return "", ErrInvalidVerificationCode
	}

	if time.Now().Unix() > pending.ExpiresAt {
		_ = s.store.Delete(pendingRegistrationKey(email))
		return "", ErrInvalidVerificationCode
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pending.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Login:    pending.Login,
		Email:    pending.Email,
		Password: string(hash),
		Role:     "user",
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	_ = s.store.Delete(pendingRegistrationKey(email))

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if user.Password == "" {
		return "", ErrGoogleOnlyAuth
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func pendingRegistrationKey(email string) string {
	return "pending_registration:" + email
}
