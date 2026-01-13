package tui

import (
	"fmt"
	"strings"

	"github.com/ManoloEsS/burrow/internal/service"
)

func (tui *Tui) handleStartServer() {
	serverPath := tui.Components.ServerPath.GetText()
	if serverPath == "" {
		serverPath = "."
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ServerStatus.SetText("[yellow]Starting server...[-]")
	})

	err := tui.Services.ServerService.StartServer(serverPath)
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to start server: %s[-]", err.Error()))
		})
		return
	}
}

func (tui *Tui) handleStopServer() {
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ServerStatus.SetText("[yellow]Stopping server...[-]")
	})

	err := tui.Services.ServerService.StopServer()
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to stop server: %s[-]", err.Error()))
		})
		return
	}
}

func (tui *Tui) updateOnServerStatusChange(status service.ServerStatus) {
	tui.Ui.QueueUpdateDraw(func() {
		tui.updateServerStatus(status)
	})
}

func (tui *Tui) updateServerStatus(status service.ServerStatus) {
	tui.State.CurrentServer = status
	var statusText string
	if status.Running {
		if strings.Contains(strings.ToLower(status.Status), "crashed") {
			statusText = fmt.Sprintf("[red]%s[-]", status.Status)
		} else {
			statusText = fmt.Sprintf("[green]%s[-]", status.Status)
		}
	} else {
		statusText = fmt.Sprintf("[red]%s[-]", status.Status)
	}
	tui.Components.ServerStatus.SetText(statusText)
}
