package screens

import (
	"strings"

	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui"
)

type MainMenu struct{}

func NewMainMenu() ui.Screen {
	return &MainMenu{}
}

func (m *MainMenu) ID() string {
	return "main"
}

func (m *MainMenu) Handle(u ui.UI, st *state.State) (string, bool) {
	line, err := u.ReadLine()
	if err != nil {
		return "main", true
	}

	cmd := strings.TrimSpace(line)
	switch cmd {
	case "create":
		u.Println("Implement create screen with make and save options")
		return "request", false
	case "retrieve":
		u.Println("Implement retrieve screen")
		return "main", false
	case "list":
		u.Println("Implement list screen")
		return "main", false
	case "quit":
		u.Println("Exiting...")
		return "main", true
	default:
		u.Printf("You typed: %s\n", cmd)
		u.Println("Commands: back, quit, create, list, retrieve")
		return "main", false
	}
}

func (m *MainMenu) Render(u ui.UI, st *state.State) {
	u.Println("Burrow")
	u.Println("This is a main screen")
	u.Println("")
	u.Println("Available commands:")
	u.Println("  create   - create new request")
	u.Println("  list     - list saved requests")
	u.Println("  retrieve - retrieve saved request")
	u.Println("  quit     - Exit application")
	u.Println("")
	u.Printf("> ")
}
