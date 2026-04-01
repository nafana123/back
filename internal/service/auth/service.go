package auth

import (
	authdto "back/internal/dto/auth"
	"back/internal/model"
	"fmt"
	"time"

	userRepository "back/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func TelegramAuth(data, botToken, jwtSecret string) (*authdto.AuthResponse, error) {
	if err := initdata.Validate(data, botToken, 24*time.Hour); err != nil {
		return nil, fmt.Errorf("Невалидный initData: %w", err)
	}

	parsed, err := initdata.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения данных пользователя: %w", err)
	}

	tgUser := parsed.User

	user := &model.User{
		ID:           tgUser.ID,
		FirstName:    tgUser.FirstName,
		LastName:     tgUser.LastName,
		Username:     tgUser.Username,
		LanguageCode: tgUser.LanguageCode,
		PhotoURL:     tgUser.PhotoURL,
		Role:         "user",
	}

	if err := userRepository.UpsertUser(user); err != nil {
		return nil, fmt.Errorf("Ошибка сохранения пользователя: %w", err)
	}

	token, err := generateJWT(user, jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("Ошибка генерации токена: %w", err)
	}

	return &authdto.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func generateJWT(user *model.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
