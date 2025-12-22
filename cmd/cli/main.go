package main

import (
	"log"
	"os"

	"github.com/ManoloEsS/burrow/internal/cli"
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/joho/godotenv"
)

func main() {
	// //connect to database
	// db, err := sql.Open("sqlite3", dbFile)
	// if err != nil {
	// 	log.Fatalf("could not connect to database: %v", err)
	// }
	// defer db.Close()
	//
	godotenv.Load()

	defaultPort := os.Getenv("DEFAULT_PORT")
	if defaultPort == "" {
		log.Fatalf("default port must be set")
	}

	cfg := config.Config{
		DefaultPort: defaultPort,
	}

	st := state.State{
		Cfg:    &cfg,
		Screen: state.MainScreen,
	}

	console := cli.NewConsole(os.Stdin, os.Stdout)
	a := app.New(console, st)

	if err := a.Run(); err != nil {
		os.Exit(1)
	}

	// for {
	// 	switch state.Screen {
	// 	case state.MainScreen:
	// 		state.ScreenState = state.HandleMainScreen(scanner)
	// 	case state.NewReqScreen:
	// 		state.ScreenState = state.HandleNewReqScreen(scanner)
	// 	case state.RetrieveScreen:
	// 		state.ScreenState = state.HandleRetrieveScreen(scanner)
	// 	}

}
