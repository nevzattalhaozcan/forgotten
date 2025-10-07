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
	Redis RedisConfig
	BookAPIs BookAPIsConfig
}

type BookAPIsConfig struct {
	GoogleBooksAPIKey string
	ISBNDBAPIKey      string
	PreferredSource   string // "google", "isbndb", "openlibrary"
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
	AutoMigrate bool
	RunMigrations bool
	MigrationsPath string
}

type JWTConfig struct {
	Secret string
	ExpirationHours int
}

type AppConfig struct {
	Name string
	Version string
}

type RedisConfig struct {
	Enabled bool
	Addr string
	Password string
	DB int
	TLS bool
	CacheTTLSeconds int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found or error loading it: %v", err)
	}

	env := getEnv("SERVER_ENVIRONMENT", "development")
    defaultAutoMigrate := env != "production"

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("PORT", "8080"),
			Environment: getEnv("SERVER_ENVIRONMENT", env),
		},
		Database: DatabaseConfig{
			URL: getEnv("DB_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			AutoMigrate:    getEnvAsBool("DB_AUTO_MIGRATE", defaultAutoMigrate),
            RunMigrations:  getEnvAsBool("DB_RUN_MIGRATIONS", false),
            MigrationsPath: getEnv("MIGRATIONS_PATH", "internal/database/migrations"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default_secret"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		App: AppConfig{
			Name: getEnv("APP_NAME", "Forgotten"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Redis: RedisConfig{
			Enabled: getEnvAsBool("REDIS_ENABLED", false),
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB: getEnvAsInt("REDIS_DB", 0),
			TLS: getEnvAsBool("REDIS_TLS", false),
			CacheTTLSeconds: getEnvAsInt("REDIS_CACHE_TTL_SECONDS", 600),
		},
		BookAPIs: BookAPIsConfig{
			GoogleBooksAPIKey: getEnv("GOOGLE_BOOKS_API_KEY", ""),
			ISBNDBAPIKey:      getEnv("ISBNDB_API_KEY", ""),
			PreferredSource:   getEnv("BOOK_API_SOURCE", "google"),
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

func getEnvAsBool(name string, defaultVal bool) bool {
	if value := os.Getenv(name); value != "" {
		switch value {
			case "1", "true", "TRUE", "True", "yes", "YES", "Yes", "on", "ON", "On":
				return true
			case "0", "false", "FALSE", "False", "no", "NO", "No", "off", "OFF", "Off":
				return false
			default:
				log.Printf("invalid boolean value for %s: %s. using default: %t", name, value, defaultVal)
		}
	}
	return defaultVal
}