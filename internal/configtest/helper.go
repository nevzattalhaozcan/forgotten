package configtest

import "github.com/nevzattalhaozcan/forgotten/internal/config"

func New() *config.Config {
    return &config.Config{
        Server: config.ServerConfig{
            Port:        "8080",
            Host:        "localhost",
            Environment: "development",
        },
        Database: config.DatabaseConfig{
            URL:          "postgres://test:test@localhost:5432/test_db?sslmode=disable",
            MaxOpenConns: 5,
            MaxIdleConns: 2,
        },
        JWT: config.JWTConfig{
            Secret:          "testsecret",
            ExpirationHours: 24,
        },
        App: config.AppConfig{
            Name:    "TestApp",
            Version: "1.0.0",
        },
    }
}