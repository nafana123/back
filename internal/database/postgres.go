package database

import (
	"fmt"

	"back/internal/config"
	"back/internal/model"

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
	if err := db.AutoMigrate(&model.User{}, &model.TgUser{}); err != nil {
		logger.Fatal("Ошибка миграции", zap.Error(err))
	}

	logger.Info("База данных подключена")
	return db
}
