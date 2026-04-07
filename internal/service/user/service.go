package user

import (
	"errors"

	"back/internal/model"
	jwtpkg "back/pkg/jwt"
	"back/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService interface {
	Register(login, email, password string) (string, error)
	Login(email, password string) (string, error)
}

type userService struct {
	userRepo   repository.UserRepository
	jwtService *jwtpkg.Service
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtpkg.NewService(jwtSecret),
	}
}

func (s *userService) Register(login, email, password string) (string, error) {
	existing, err := s.userRepo.GetByLoginOrEmail(login, email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		if existing.Login == login {
			return "", ErrLoginAlreadyExists
		}
		if existing.Email == email {
			return "", ErrEmailAlreadyExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Login:    login,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) Login(email, password string) (string, error) {
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
