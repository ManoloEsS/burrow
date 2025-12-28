package screens

import (
	"strings"

	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui"
)

type RequestScreen struct{}

func NewRequestScreen() ui.Screen {
	return &RequestScreen{}
}

func (rs *RequestScreen) ID() string {
	return "request"
}

func (rs *RequestScreen) Handle(u ui.UI, st *state.State) (string, bool) {
	line, err := u.ReadLine()
	if err != nil {
		return "request", true
	}

	cmd := strings.TrimSpace(line)
	switch cmd {
	case "back":
		return "main", false
	case "quit":
		return "request", true
	default:
		u.Printf("You typed: %s\n", cmd)
		u.Println("Commands: back, quit")
		return "request", false
	}
}

func (rs *RequestScreen) Render(u ui.UI, st *state.State) {
	u.Println("Request Screen")
	u.Println("Implement request input")
	u.Println()
}
