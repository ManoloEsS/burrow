package tui

import (
	"fmt"

	"github.com/ManoloEsS/burrow/internal/service"
)

func (tui *Tui) handleStartServer() {
	serverPath := tui.Components.ServerPath.GetText()
	if serverPath == "" {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ServerStatus.SetText("[red]Filepath is empty...[-]")
		})
		return
	}
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ServerStatus.SetText("starting server")
	})

	err := tui.ServerService.StartServer(serverPath, tui.Config.App.DefaultPort, tui.ServerUpdateChannel)
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to start server: %s[-]", err.Error()))
		})
		return
	}
}

func (tui *Tui) handleStopServer() {
	err := tui.ServerService.StopServer()
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to stop server: %s[-]", err.Error()))
		})
		return
	}
}

func (tui *Tui) handleServerEvent(event service.UIEvent) {
	tui.Ui.QueueUpdateDraw(func() {
		switch event.Type {
		case "error":
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[red]%s[-]", event.Message))
		case "update":
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[green]%s[-]", event.Message))
		default:
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[yellow]%s[-]", event.Message))
		}
	})
}

func (tui *Tui) serverUpdateListener() {
	for event := range tui.ServerUpdateChannel {
		tui.handleServerEvent(event)
	}
}
