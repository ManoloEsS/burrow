package app

import (
	"context"

	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui"
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

	// TODO: REPL implementation with screens
}
