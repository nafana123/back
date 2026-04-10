package database

import (
	"back/internal/model"
	"fmt"

	"back/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(logger *zap.Logger, cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Ошибка подключения к базе", zap.Error(err))
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(&model.User{}, &model.TgUser{}, &model.Game{}, &model.Tournament{}, &model.Participant{}); err != nil {
		logger.Fatal("Ошибка миграции", zap.Error(err))
	}

	logger.Info("База данных подключена")
	return db
}
