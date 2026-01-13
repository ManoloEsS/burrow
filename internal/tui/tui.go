package tui

import (
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

type Tui struct {
	Ui         *tview.Application
	Services   *service.Services
	Components *UIComponents
	State      *UIState
}

func NewTui() *Tui {
	return &Tui{
		Ui:    tview.NewApplication(),
		State: &UIState{},
	}
}

func (tui *Tui) Initialize() error {
	tui.Components = createTuiLayout()
	tui.setupKeybindings()
	tui.loadSavedRequests()
	tui.updateServerStatus(tui.Services.ServerService.GetStatus())

	tui.focusForm()

	return nil
}

func (tui *Tui) Start() error {
	return tui.Ui.SetRoot(tui.Components.MainLayout, true).EnableMouse(true).Run()
}
