package tui

import (
	"fmt"
	"strings"

	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI wraps the tview application and components for input capture, the services layer, and the UIState. It
// also keeps track of the focused subcomponents in the form component, as well as whether the form or list
// are focused to manage keybindings exclusive to each
type UI struct {
	app        *tview.Application
	services   *service.Services
	components *UIComponents
	state      *UIState

	// Focus tracking fields
	currentFormFocusIndex int
	formFocused           bool
	listFocused           bool
}

// Creates a new UI struct
func NewUI(services *service.Services) *UI {
	return &UI{
		app:                   tview.NewApplication(),
		services:              services,
		state:                 &UIState{},
		currentFormFocusIndex: 0,
		formFocused:           false,
		listFocused:           false,
	}
}

// Initializes the tui by creating the layout, setting up service callbacks for events, keybindings.
// It loads and renders saved requests on start, as well as the default server status.
func (ui *UI) Initialize() error {
	ui.components = createTuiLayout()
	ui.setupServiceCallbacks()
	ui.setupKeybindings()
	ui.loadSavedRequests()
	ui.updateServerStatus(ui.services.ServerService.GetStatus())

	ui.focusForm()

	return nil
}

// Sets up update callbacks for each service to update the UI after a service call.
func (ui *UI) setupServiceCallbacks() {
	ui.services.RequestService.SetUpdateCallback(ui.onResponseReceived)
	ui.services.ServerService.SetUpdateCallback(ui.onServerStatusChanged)
}

// Sets focus on form component to make keybindings specific to form avialable.
// Also keeps track of subcomponent in focus when shifting focus out and back into the form component.
func (ui *UI) focusForm() {
	ui.formFocused = true
	ui.listFocused = false
	ui.focusSpecificFormComponent(ui.currentFormFocusIndex)
}

// Sets focus on server status component input.
func (ui *UI) focusServerInput() {
	ui.formFocused = false
	ui.listFocused = false
	ui.app.SetFocus(ui.components.ServerPath)
}

// Sets focus on request name input for saving request to db.
func (ui *UI) focusRequestNameInput() {
	ui.formFocused = false
	ui.listFocused = false
	ui.app.SetFocus(ui.components.NameInput)
}

// Sets focus on saved requests list to make specific keybindings available.
func (ui *UI) focusRequestList() {
	ui.formFocused = false
	ui.listFocused = true
	ui.app.SetFocus(ui.components.RequestList)
}

// Enables up and down navigation of subcomponents in form component.
func (ui *UI) navigateForm(forward bool) {
	subcompCount := ui.components.Form.GetFormItemCount()
	if forward {
		ui.currentFormFocusIndex = (ui.currentFormFocusIndex + 1) % subcompCount
	} else {
		ui.currentFormFocusIndex = (ui.currentFormFocusIndex - 1 + 6) % subcompCount
	}

	ui.focusSpecificFormComponent(ui.currentFormFocusIndex)
}

// Sets focus on specific subcomponent in form component when navigating through subcomponents.
func (ui *UI) focusSpecificFormComponent(index int) {
	component := ui.components.Form.GetFormItem(index)
	if component != nil {
		ui.app.SetFocus(component)
	}
}

// Allows for up and down navigation in list component
func (ui *UI) navigateList(direction int) {
	currentList := ui.components.RequestList
	currentIndex := currentList.GetCurrentItem()
	itemCount := currentList.GetItemCount()

	if itemCount == 0 {
		return
	}

	newItem := (currentIndex + direction + itemCount) % itemCount
	currentList.SetCurrentItem(newItem)
}

// Loads saved requests and renders them in saved requests component
func (ui *UI) loadSavedRequests() {
	err := ui.services.RequestService.GetSavedRequests()
	if err != nil {
		return
	}

	ui.state.RequestHistory = nil

	// // Load requests into list
	// for _, req := range requests {
	// 	itemText := fmt.Sprintf("%s %s", req.Method, req.URL)
	// 	secondaryText := req.URL
	// 	ui.components.RequestList.AddItem(itemText, secondaryText, 0, nil)
	// }
}

// Callback function to render received response to response viewer
func (ui *UI) onResponseReceived(response *service.Response) {
	ui.app.QueueUpdateDraw(func() {
		ui.state.CurrentResponse = response
		responseText := fmt.Sprintf(
			"[green]Latest Response[%d: %d]\n\n%s[-]",
			response.RequestID,
			response.StatusCode,
			response.Body,
		)
		ui.components.ResponseView.SetText(responseText)
	})
}

// Callback function to render changes in server status
func (ui *UI) onServerStatusChanged(status service.ServerStatus) {
	ui.app.QueueUpdateDraw(func() {
		ui.updateServerStatus(status)
	})
}

// Updates server status in ui state and renders text in server status text view
func (ui *UI) updateServerStatus(status service.ServerStatus) {
	ui.state.CurrentServer = status
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
	ui.components.ServerStatus.SetText(statusText)
}

// Main keybindings configuration for tui, navigation keybindings update focus, while event keybindings
// call handlers that perform a behavior and callbacks to update tui
func (ui *UI) setupKeybindings() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Events
		case tcell.KeyCtrlS:
			go ui.handleSendRequest()
			return nil
		case tcell.KeyCtrlQ:
			ui.app.Stop()
			return nil
		case tcell.KeyF5:
			go ui.handleStartServer()
			return nil
		case tcell.KeyF6:
			go ui.handleStopServer()
			return nil

		// Navigation
		case tcell.KeyCtrlF:
			ui.focusForm()
			return nil
		case tcell.KeyCtrlG:
			ui.focusServerInput()
			return nil
		case tcell.KeyCtrlL:
			ui.focusRequestList()
			return nil
		case tcell.KeyCtrlT:
			ui.focusRequestNameInput()
			return nil
		case tcell.KeyCtrlN:
			if ui.formFocused {
				ui.navigateForm(true)
			}
			return nil
		case tcell.KeyCtrlP:
			if ui.formFocused {
				ui.navigateForm(false)
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				if ui.listFocused {
					ui.navigateList(1)
					return nil
				}
				return event
			case 'k':
				if ui.listFocused {
					ui.navigateList(-1)
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
func (ui *UI) handleSendRequest() {
	req := ui.getCurrentRequest()

	ui.app.QueueUpdateDraw(func() {
		ui.components.ResponseView.SetText("[yellow]Sending request...[-]")
	})

	_, err := ui.services.RequestService.SendRequest(req)
	if err != nil {
		ui.app.QueueUpdateDraw(func() {
			ui.components.ResponseView.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
		})
		return
	}

	ui.app.QueueUpdateDraw(func() {
		itemText := fmt.Sprintf("%s %s", req.Method, req.URL)
		secondaryText := req.URL
		ui.components.RequestList.AddItem(itemText, secondaryText, 0, nil)
	})

	ui.updateRequestHistory()
}

// Handler to start a go server on specified path and update server status.
func (ui *UI) handleStartServer() {
	serverPath := ui.components.ServerPath.GetText()
	if serverPath == "" {
		serverPath = "localhost:8080"
	}

	ui.app.QueueUpdateDraw(func() {
		ui.components.ServerStatus.SetText("[yellow]Starting server...[-]")
	})

	err := ui.services.ServerService.StartServer(serverPath)
	if err != nil {
		ui.app.QueueUpdateDraw(func() {
			ui.components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to start server: %s[-]", err.Error()))
		})
		return
	}
}

// Handler to stop running server
func (ui *UI) handleStopServer() {
	ui.app.QueueUpdateDraw(func() {
		ui.components.ServerStatus.SetText("[yellow]Stopping server...[-]")
	})

	err := ui.services.ServerService.StopServer()
	if err != nil {
		ui.app.QueueUpdateDraw(func() {
			ui.components.ServerStatus.SetText(fmt.Sprintf("[red]Failed to stop server: %s[-]", err.Error()))
		})
		return
	}
}

// Captures inputs from form to create a new request ready to be saved or sent
func (ui *UI) getCurrentRequest() *domain.Request {
	_, method := ui.components.MethodDropdown.GetCurrentOption()

	url := ui.components.URLInput.GetText()

	headersText := ui.components.HeadersText.GetText()

	paramsText := ui.components.ParamsText.GetText()

	_, bodyType := ui.components.BodyType.GetCurrentOption()

	body := ui.components.BodyText.GetText()

	newRequest := domain.Request{}

	err := newRequest.ParseMethod(method)
	if err != nil {
		return &newRequest
	}

	// Create basic config for URL parsing - this is a temporary solution

	err = newRequest.ParseUrl(ui.services.Config, url)
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
func (ui *UI) updateRequestHistory() {
	err := ui.services.RequestService.GetSavedRequests()
	if err != nil {
		return
	}

	ui.state.RequestHistory = nil

}

// Starts app setting layout as root and enables mouse interactions
func (ui *UI) Start() error {
	return ui.app.SetRoot(ui.components.MainLayout, true).EnableMouse(true).Run()
}
