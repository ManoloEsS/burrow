package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// load and validate config
	cfg := config.LoadFromEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// initialize database
	databaseInstance, err := database.NewDatabase(cfg.DbPath, cfg.DbString, cfg.DbMigrationsDir)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer databaseInstance.Close()

	setupGracefulShutdown(databaseInstance)

	// initialize ui and run app
}

func setupGracefulShutdown(db *database.Database) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		db.Close()
		os.Exit(0)
	}()
}
