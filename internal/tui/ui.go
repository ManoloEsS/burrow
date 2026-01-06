package tui

import (
	"fmt"
	"strings"

	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/gdamore/tcell/v2"
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

// Sets focus on form component to make keybindings specific to form avialable.
// Also keeps track of subcomponent in focus when shifting focus out and back into the form component.
func (tui *Tui) focusForm() {
	tui.State.CurrentFocused = tui.Components.Form
	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
}

// Sets focus on server status component input.
func (tui *Tui) focusServerInput() {
	tui.State.CurrentFocused = tui.Components.ServerPath
	tui.Ui.SetFocus(tui.Components.ServerPath)
}

// Sets focus on request name input for saving request to db.
func (tui *Tui) focusRequestNameInput() {
	tui.State.CurrentFocused = tui.Components.NameInput
	tui.Ui.SetFocus(tui.Components.NameInput)
}

// Sets focus on saved requests list to make specific keybindings available.
func (tui *Tui) focusRequestList() {
	tui.State.CurrentFocused = tui.Components.RequestList
	tui.Ui.SetFocus(tui.Components.RequestList)
}

// Enables up and down navigation of subcomponents in form component.
func (tui *Tui) navigateForm(forward bool) {
	subcompCount := tui.Components.Form.GetFormItemCount()
	if forward {
		tui.State.CurrentFormFocusIndex = (tui.State.CurrentFormFocusIndex + 1) % subcompCount
	} else {
		tui.State.CurrentFormFocusIndex = (tui.State.CurrentFormFocusIndex - 1 + 6) % subcompCount
	}

	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
}

// Sets focus on specific subcomponent in form component when navigating through subcomponents.
func (tui *Tui) focusSpecificFormComponent(index int) {
	component := tui.Components.Form.GetFormItem(index)
	if component != nil {
		tui.Ui.SetFocus(component)
	}
}

// Allows for up and down navigation in list component
func (tui *Tui) navigateList(direction int) {
	currentList := tui.Components.RequestList
	currentIndex := currentList.GetCurrentItem()
	itemCount := currentList.GetItemCount()

	if itemCount == 0 {
		return
	}

	newItem := (currentIndex + direction + itemCount) % itemCount
	currentList.SetCurrentItem(newItem)
}

// Loads saved requests and renders them in saved requests component
func (tui *Tui) loadSavedRequests() {
	err := tui.Services.RequestService.GetSavedRequests()
	if err != nil {
		return
	}

	tui.State.RequestHistory = nil

	// // Load requests into list
	// for _, req := range requests {
	// 	itemText := fmt.Sprintf("%s %s", req.Method, req.URL)
	// 	secondaryText := req.URL
	// 	a.Components.RequestList.AddItem(itemText, secondaryText, 0, nil)
	// }
}

// Callback function to render received response to response viewer
func (tui *Tui) UpdateOnReceiveResponse(response *domain.Response) {
	tui.Ui.QueueUpdateDraw(func() {
		tui.State.CurrentResponse = response
		responseText := fmt.Sprintf(
			"[green]Latest Response[%d: %d]\n\n%s[-]",
			response.RequestID,
			response.StatusCode,
			response.Body,
		)
		tui.Components.ResponseView.SetText(responseText)
	})
}

func (tui *Tui) UpdateOnServerStatusChange(status service.ServerStatus) {
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

// Main keybindings configuration for tui, navigation keybindings update focus, while event keybindings
// call handlers that perform a behavior and callbacks to update tui
func (tui *Tui) setupKeybindings() {
	tui.Ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Events
		case tcell.KeyCtrlS:
			go tui.handleSendRequest()
			return nil
		case tcell.KeyCtrlQ:
			tui.Ui.Stop()
			return nil
		case tcell.KeyF5:
			go tui.handleStartServer()
			return nil
		case tcell.KeyF6:
			go tui.handleStopServer()
			return nil

		// Navigation
		case tcell.KeyCtrlF:
			tui.focusForm()
			return nil
		case tcell.KeyCtrlG:
			tui.focusServerInput()
			return nil
		case tcell.KeyCtrlL:
			tui.focusRequestList()
			return nil
		case tcell.KeyCtrlT:
			tui.focusRequestNameInput()
			return nil
		case tcell.KeyCtrlN:
			if tui.State.CurrentFocused == tui.Components.Form {
				tui.navigateForm(true)
			}
			return nil
		case tcell.KeyCtrlP:
			if tui.State.CurrentFocused == tui.Components.Form {
				tui.navigateForm(false)
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				if tui.State.CurrentFocused == tui.Components.RequestList {
					tui.navigateList(1)
					return nil
				}
				return event
			case 'k':
				if tui.State.CurrentFocused == tui.Components.RequestList {
					tui.navigateList(-1)
					return nil
				}
				return event
			default:
				return event
			}
		default:
			return event
		}
	})
}

// Handler to send http requests from captured tui data.
func (tui *Tui) handleSendRequest() {
	req := tui.getCurrentRequest()

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ResponseView.SetText("[yellow]Sending request...[-]")
	})

	_, err := tui.Services.RequestService.SendRequest(req)
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ResponseView.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
		})
		return
	}

	tui.Ui.QueueUpdateDraw(func() {
		itemText := fmt.Sprintf("%s %s", req.Method, req.URL)
		secondaryText := req.URL
		tui.Components.RequestList.AddItem(itemText, secondaryText, 0, nil)
	})

	tui.updateRequestHistory()
}

// Handler to start a go server on specified path and update server status.
func (tui *Tui) handleStartServer() {
	serverPath := tui.Components.ServerPath.GetText()
	if serverPath == "" {
		serverPath = "localhost:8080"
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

// Handler to stop running server
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

// Captures inputs from form to create a new request ready to be saved or sent
func (tui *Tui) getCurrentRequest() *domain.Request {
	_, method := tui.Components.MethodDropdown.GetCurrentOption()

	url := tui.Components.URLInput.GetText()

	headersText := tui.Components.HeadersText.GetText()

	paramsText := tui.Components.ParamsText.GetText()

	_, bodyType := tui.Components.BodyType.GetCurrentOption()

	body := tui.Components.BodyText.GetText()

	newRequest := domain.Request{}

	err := newRequest.ParseMethod(method)
	if err != nil {
		return &newRequest
	}

	// Create basic config for URL parsing - this is a temporary solution

	err = newRequest.ParseUrl(tui.Services.Config, url)
	if err != nil {
		return &newRequest
	}

	err = newRequest.ParseHeaders(headersText)
	if err != nil {
		return &newRequest
	}

	err = newRequest.ParseParams(paramsText)
	if err != nil {
		return &newRequest
	}

	err = newRequest.ParseBodyType(bodyType)
	if err != nil {
		return &newRequest
	}

	err = newRequest.ParseBody(body)
	if err != nil {
		return &newRequest
	}

	return &newRequest
}

// Updates saved requests list
func (tui *Tui) updateRequestHistory() {
	err := tui.Services.RequestService.GetSavedRequests()
	if err != nil {
		return
	}

	tui.State.RequestHistory = nil

}

// Starts app setting layout as root and enables mouse interactions
func (tui *Tui) Start() error {
	return tui.Ui.SetRoot(tui.Components.MainLayout, true).EnableMouse(true).Run()
}
