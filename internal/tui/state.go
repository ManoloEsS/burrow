package tui

import (
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

type UIState struct {
	CurrentRequest        *domain.Request
	CurrentServer         service.ServerStatus
	RequestHistory        []*domain.Request
	CurrentResponse       *domain.Response
	CurrentFormFocusIndex int
	CurrentFocused        tview.Primitive
}
