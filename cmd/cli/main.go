package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ManoloEsS/burrow/internal/app"
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui/console"
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

	st := &state.State{
		Cfg: cfg,
		DB:  databaseInstance,
		Screen: state.ScreenState{
			CurrentScreen: "main",
			SelectedID:    "",
		},
		Requests: make(map[string]domain.Request),
	}

	// initialize ui and run app
	consoleUI := console.NewConsole(os.Stdin, os.Stdout)
	a := app.New(consoleUI, st)

	if err := a.Run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
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
