package main

import (
	"fmt"
)

func main() {
	// //connect to database
	// db, err := sql.Open("sqlite3", dbFile)
	// if err != nil {
	// 	log.Fatalf("could not connect to database: %v", err)
	// }
	// defer db.Close()
	//

	// cfg := config.Config{
	// 	DefaultPort: DefaultPort,
	// }

	// state := &state.State{
	// 	Cfg:    &cfg,
	// 	Screen: state.MainScreen,
	// }
	//
	// scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to burrow")

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
