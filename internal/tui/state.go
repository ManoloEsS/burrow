package tui

import (
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/rivo/tview"
)

type UIState struct {
	CurrentRequest        *domain.Request
	SavedRequests         []*domain.Request
	CurrentResponse       *domain.Response
	CurrentFormFocusIndex int
	CurrentFocused        tview.Primitive
}
