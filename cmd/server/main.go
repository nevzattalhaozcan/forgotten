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