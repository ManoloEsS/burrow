package tui

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

type Tui struct {
	Ui            *tview.Application
	HttpService   service.HttpClientService
	ServerService service.ServerService
	Components    *UIComponents
	State         *UIState
	Config        *config.Config
}

func NewTui(cfg *config.Config) *Tui {
	return &Tui{
		Ui:     tview.NewApplication(),
		State:  &UIState{},
		Config: cfg,
	}

}

func (tui *Tui) Initialize() error {
	tui.Components = createTuiLayout()
	tui.setupKeybindings()
	tui.loadSavedRequests()
	tui.updateServerStatus(tui.ServerService.HealthCheck())

	tui.focusForm()

	return nil
}

func (tui *Tui) Start() error {
	return tui.Ui.SetRoot(tui.Components.MainLayout, true).EnableMouse(true).Run()
}
