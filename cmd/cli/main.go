package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/ManoloEsS/burrow/internal/tui"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := config.LoadFromEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := database.NewDatabase(cfg.DbPath, cfg.DbString, cfg.DbMigrationsDir)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	// Create empty UI to get callbacks
	ui := tui.NewTui()

	// Create services with UI callbacks
	services := service.NewServices(db, cfg,
		ui.UpdateOnReceiveResponse,
		ui.UpdateOnServerStatusChange)

	// Inject services into UI
	ui.Services = services

	// Initialize UI
	if err := ui.Initialize(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}

	// Setup graceful shutdown
	setupShutdown(db)

	// Start the application
	if err := ui.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func setupShutdown(db *database.Database) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		db.Close()
		os.Exit(0)
	}()
}
