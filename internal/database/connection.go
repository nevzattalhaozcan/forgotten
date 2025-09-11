package database

import (
	"log"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error)  {
	logLevel := logger.Silent
	if cfg.Server.Environment == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	if err := autoMigrate(db); err != nil {
		log.Printf("migration error: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
	)
}