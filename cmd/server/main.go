// @title Forgotten API
// @version 1.0
// @description A Twitter-like social media API built with Go
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @termsOfService http://swagger.io/terms/

package main

import (
	"fmt"
	"log"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/database"
	"github.com/nevzattalhaozcan/forgotten/internal/handlers"
	"github.com/nevzattalhaozcan/forgotten/pkg/logger"
)

func main() {
	cfg := config.Load()

	if err := logger.Init(cfg.Server.Environment); err != nil {
        log.Fatalf("failed to initialize logger: %v", err)
    }
    defer logger.Sync()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	server := handlers.NewServer(db, cfg)
	

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("starting server on %s", addr)

	if err := server.Start(addr); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}