package tui

import (
	"log"
	"os"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

type Tui struct {
	Ui                  *tview.Application
	HttpService         service.HttpClientService
	ServerService       service.ServerService
	Components          *UIComponents
	State               *UIState
	Config              *config.Config
	logger              *log.Logger
	logFile             *os.File
	ServerUpdateChannel chan service.UIEvent
}

func NewTui(cfg *config.Config) *Tui {
	logFile, err := os.Create(cfg.Paths.LogFile)
	if err != nil {
		return &Tui{
			Ui:                  tview.NewApplication(),
			State:               &UIState{},
			Config:              cfg,
			logger:              nil,
			ServerUpdateChannel: make(chan service.UIEvent, 30),
		}
	}

	return &Tui{
		Ui:                  tview.NewApplication(),
		State:               &UIState{},
		Config:              cfg,
		logger:              log.New(logFile, "[TUI] ", log.LstdFlags),
		logFile:             logFile,
		ServerUpdateChannel: make(chan service.UIEvent, 30),
	}
}

func (tui *Tui) Initialize() error {
	tui.Components = createTuiLayout(tui.Config)
	tui.setupKeybindings()
	tui.loadSavedRequests()
	tui.focusForm()
	go tui.serverUpdateListener()

	return nil
}

func (tui *Tui) Start() error {
	return tui.Ui.SetRoot(tui.Components.Pages, true).EnableMouse(true).Run()
}

func (tui *Tui) Stop() {
	close(tui.ServerUpdateChannel)
	tui.Ui.Stop()
}

func (tui *Tui) Close() error {
	if tui.logFile != nil {
		return tui.logFile.Close()
	}
	return nil
}
