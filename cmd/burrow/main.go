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
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := database.NewDatabase(cfg.Database.Path, cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	ui := tui.NewTui(cfg)

	ui.HttpService = service.NewHttpClientService(db)
	ui.ServerService = service.NewServerService()

	if err := ui.Initialize(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}

	setupShutdown(db)

	if err := ui.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func setupShutdown(db *database.Database) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down database safely")
		db.Close()
		os.Exit(0)
	}()
}
