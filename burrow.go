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
	defer func() { _ = db.Close() }()

	ui := tui.NewTui(cfg)
	defer ui.Close()

	ui.HttpService = service.NewHttpClientService(db)
	ui.ServerService = service.NewServerService()

	if err := ui.Initialize(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}

	setupShutdown()

	if err := ui.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func setupShutdown() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutdown signal received, exiting...")
		os.Exit(0)
	}()
}
