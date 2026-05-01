package google

import (
	"back/internal/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	authdto "back/internal/dto/auth"
	"back/internal/model"
	userrepo "back/internal/repository/user"
	jwtpkg "back/pkg/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type GoogleService struct {
	cfg         *config.Config
	oauthConfig *oauth2.Config
	userRepo    *userrepo.UserRepository
	jwtService  *jwtpkg.Service
}

func NewGoogleService(cfg *config.Config, userRepo *userrepo.UserRepository) *GoogleService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleCallbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleService{
		cfg:         cfg,
		oauthConfig: oauthConfig,
		userRepo:    userRepo,
		jwtService:  jwtpkg.NewService(cfg.JWTSecret),
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *GoogleService) GoogleValidate(ctx context.Context, code string) (*authdto.AuthResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("ошибка получения google auth code")
	}

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("ошибка обмена google token: %w", err)
	}

	client := s.oauthConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo вернул неуспешный статус: %d", resp.StatusCode)
	}

	var user authdto.GoogleUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("ошибка декодирования google userinfo: %w", err)
	}

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Name = strings.TrimSpace(user.Name)

	if user.Email == "" {
		return nil, errors.New("google email не найден")
	}

	if !user.VerifiedEmail {
		return nil, errors.New("google email не подтверждён")
	}

	login := user.Name
	if login == "" {
		login = strings.Split(user.Email, "@")[0]
	}

	dbUser, err := s.userRepo.GetByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		dbUser = &model.User{
			Login:    login,
			Email:    user.Email,
			Password: "",
			PhotoURL: user.Picture,
			Role:     "user",
		}

		if err := s.userRepo.CreateUser(dbUser); err != nil {
			return nil, err
		}
	} 

	jwtToken, err := s.jwtService.GenerateToken(dbUser.ID, dbUser.Role)
	if err != nil {
		return nil, err
	}

	return &authdto.AuthResponse{
		Token: jwtToken,
	}, nil
}
