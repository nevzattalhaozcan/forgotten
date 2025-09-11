package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	JWT JWTConfig
	App AppConfig
}

type ServerConfig struct {
	Port string
	Host string
	Environment string
}

type DatabaseConfig struct {
	URL string
	MaxOpenConns int
	MaxIdleConns int
}

type JWTConfig struct {
	Secret string
	ExpirationHours int
}

type AppConfig struct {
	Name string
	Version string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found or error loading it: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default_secret"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		App: AppConfig{
			Name: getEnv("APP_NAME", "Forgotten"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	if value := os.Getenv(name); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("invalid integer value for %s: %s. using default: %d", name, value, defaultVal)
	}
	return defaultVal
}