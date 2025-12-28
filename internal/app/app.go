package app

import (
	"context"

	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui"
	"github.com/ManoloEsS/burrow/internal/ui/screens"
)

type App struct {
	UI     ui.UI
	State  *state.State
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ui ui.UI, st *state.State) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		UI:     ui,
		State:  st,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() error {
	currentScreen := a.State.Screen.CurrentScreen
	shouldExit := false

	for !shouldExit {
		var screen ui.Screen

		switch currentScreen {
		case "main":
			screen = screens.NewMainMenu()
		case "request":
			screen = screens.NewRequestScreen()
		case "retrieve":
			// TODO: Implement retrieve screen
			a.UI.Println("Retrieve screen not implemented yet")
			a.State.Screen.SetScreen("main")
			continue
		default:
			a.UI.Println("Unknown screen, returning to main menu")
			a.State.Screen.SetScreen("main")
			currentScreen = "main"
			continue
		}

		screen.Render(a.UI, a.State)

		nextScreen, exit := screen.Handle(a.UI, a.State)

		a.State.Screen.SetScreen(nextScreen)
		currentScreen = nextScreen
		shouldExit = exit
	}

	return nil
}
