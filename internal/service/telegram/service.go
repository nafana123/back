package telegram

import (
	authdto "back/internal/dto/auth"
	"back/internal/model"
	"back/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type TelegramService interface {
	TelegramAuth(data, botToken, jwtSecret string) (*authdto.AuthResponse, error)
}

type telegramService struct {
	userRepo repository.TgUserRepository
}

func NewTelegramService(userRepo repository.TgUserRepository) TelegramService {
	return &telegramService{
		userRepo: userRepo,
	}
}

func (s *telegramService) TelegramAuth(data, botToken, jwtSecret string) (*authdto.AuthResponse, error) {
	if err := initdata.Validate(data, botToken, 24*time.Hour); err != nil {
		return nil, fmt.Errorf("Невалидный initData: %w", err)
	}

	parsed, err := initdata.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения данных пользователя: %w", err)
	}

	tgUser := parsed.User

	user := &model.TgUser{
		ID:           tgUser.ID,
		FirstName:    tgUser.FirstName,
		LastName:     tgUser.LastName,
		Username:     tgUser.Username,
		LanguageCode: tgUser.LanguageCode,
		PhotoURL:     tgUser.PhotoURL,
		Role:         "user",
	}

	if err := s.userRepo.UpsertUser(user); err != nil {
		return nil, fmt.Errorf("Ошибка сохранения пользователя: %w", err)
	}

	token, err := s.generateJWT(user, jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("Ошибка генерации токена: %w", err)
	}

	return &authdto.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *telegramService) generateJWT(user *model.TgUser, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}