package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/GoldenDeals/DepGit/internal/config"
	"github.com/GoldenDeals/DepGit/internal/database"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
)

var mainLogger = logger.New("main")

func main() {
	// Load configuration from .env file and environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create necessary directories
	dbDir := filepath.Dir(cfg.GetDatabasePath())
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		mainLogger.Fatalf("Failed to create database directory: %v", err)
	}

	// Create migrations directory if specified
	if cfg.DB.MigrationsPath != "" {
		if err := os.MkdirAll(cfg.DB.MigrationsPath, 0755); err != nil {
			mainLogger.Fatalf("Failed to create migrations directory: %v", err)
		}
	}

	// Initialize the database
	mainLogger.Info("Initializing database...")
	db := &database.DB{}
	if err := db.Init(cfg); err != nil {
		mainLogger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	mainLogger.Info("DepGit server starting")
	// TODO: Add server initialization and other startup logic here

	// Keep the server running
	select {}
}
