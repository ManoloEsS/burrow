package screens

import (
	"github.com/ManoloEsS/burrow/internal/state"
	"github.com/ManoloEsS/burrow/internal/ui"
)

type Screen interface {
	ID() string
	Render(u ui.UI, st *state.State)
	Handle(u ui.UI, st *state.State) (string, bool)
}
