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

func Connect(cfg *config.Config) (*gorm.DB, error) {
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

	if cfg.Server.Environment == "production" {
		if cfg.Database.RunMigrations {
			if err := RunMigrations(db, cfg.Database.MigrationsPath); err != nil {
				log.Printf("migration error: %v", err)
				return nil, err
			}
		} else {
			log.Println("Production: skipping AutoMigrate; run SQL migrations externally (make migrate-up).")
		}
	} else {
		// dev/test automigrate if enabled
		if cfg.Database.AutoMigrate {
			if err := autoMigrate(db); err != nil {
				log.Printf("autoMigrate error: %v", err)
				return nil, err
			}
		}
	}

	if cfg.Server.Environment == "development" {
        if err := SeedForTest(db); err != nil {
            log.Printf("seeding error: %v", err)
            // non-fatal
        }
    }

	log.Println("Database connected successfully")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
		models.Book{},
		models.Club{},
		models.ClubMembership{},
		models.Event{},
		models.EventRSVP{},
		models.Comment{},
		models.CommentLike{},
		models.Post{},
		models.PostLike{},
		models.UserBookProgress{},
		models.ClubBookAssignment{},
		models.ReadingLog{},
		models.ClubRating{},
	)
}
