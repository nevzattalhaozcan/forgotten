package database

import (
	"github.com/nevzattalhaozcan/forgotten/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, migrationsPath string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

func RollbackMigrations(db *gorm.DB, migrationsPath string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Steps(-steps); err != nil {
		return err
	}

	logger.Info("Database rollback completed successfully")
	return nil
}