package test_helpers

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Club{},
        &models.Event{},
        &models.EventRSVP{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CleanupTestDB(db *gorm.DB) error {
	tables := []string{"users"}
	
	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
	}
	
	return nil
}