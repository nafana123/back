package user

import (
	userdto "back/internal/dto/user"
	mail "back/internal/mailer"
	"errors"
	"fmt"
	"math/rand"
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
	subject := "Код подтверждения регистрации"
	if err := s.mailer.Send(email, subject, code); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailDeliveryFailed, err)
	}

	if err := s.store.Set(email, code); err != nil {
		return err
	}

	return nil
}

func (s *Service) CompleteVerification(body userdto.VerifyRequest) (string, error) {
	code, err := s.store.Get(body.Email)
	if err != nil || code != body.Code {
		return "", ErrInvalidVerificationCode
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Login:    body.Login,
		Email:    body.Email,
		Password: string(hash),
		Role:     "user",
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	_ = s.store.Delete(body.Email)

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
