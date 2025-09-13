package database

import (
	"log"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/pkg/metrics"
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

	// Periodically publish DB connection stats
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            stats := sqlDB.Stats()
            metrics.UpdateDBMetrics(stats.OpenConnections, stats.Idle)
        }
    }()

	if err := autoMigrate(db); err != nil {
		log.Printf("migration error: %v", err)
		return nil, err
	}

	// Seed the database with initial data for testing
	if cfg.Server.Environment == "development" {
		if err := SeedForTest(db); err != nil {
			log.Printf("seeding error: %v", err)
			// Don't fail if seeding fails, just log it
		}
	}

	log.Println("Database connected successfully")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
	)
}