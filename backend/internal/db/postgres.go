package db

import (
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/config"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgres(cfg config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return database, nil
}

func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Task{},
	)
}
