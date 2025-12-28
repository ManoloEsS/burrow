package ui

import (
	"github.com/ManoloEsS/burrow/internal/state"
)

type Screen interface {
	ID() string
	Render(u UI, st *state.State)
	Handle(u UI, st *state.State) (string, bool)
}
