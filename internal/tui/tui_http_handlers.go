package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/ManoloEsS/burrow/internal/domain"
)

func (tui *Tui) handleLoadRequest() {
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Populating...")
	})

	if tui.Components.RequestList.GetItemCount() < 1 {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.StatusText.SetText("No saved requests")
		})
		return
	}
	index := tui.Components.RequestList.GetCurrentItem()

	text, _ := tui.Components.RequestList.GetItemText(index)

	textParts := strings.Split(text, "|")

	name := strings.TrimSpace(textParts[0])

	for _, request := range tui.State.SavedRequests {
		if request.Name == name {
			tui.Ui.QueueUpdateDraw(func() {
				tui.populateRequest(request)
			})
		}
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Request loaded")
	})

}

func (tui *Tui) handleDeleteRequest() {
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Deleting...")
	})

	if tui.Components.RequestList.GetItemCount() < 1 {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.StatusText.SetText("No saved requests")
		})
		return

	}

	index := tui.Components.RequestList.GetCurrentItem()
	text, _ := tui.Components.RequestList.GetItemText(index)
	err := tui.HttpService.DeleteRequest(text)
	if err != nil {
		log.Printf("could not delete request: %v", err)
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.StatusText.SetText(fmt.Sprintf("Error %v", err))
		})
		return
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.loadSavedRequests()
	})
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Request deleted")
	})
}

func (tui *Tui) handleSaveRequest() {
	err := tui.getCurrentRequest()
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.StatusText.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
		})
		return
	}
	if tui.State.CurrentRequest.Name == "" {
		tui.State.CurrentRequest.Name = fmt.Sprintf("unnamed%d", tui.Components.RequestList.GetItemCount()+1)
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Saving request...")
	})

	err = tui.HttpService.SaveRequest(tui.State.CurrentRequest)
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.StatusText.SetText(fmt.Sprintf("Error: %s", err))
		})
		return
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.loadSavedRequests()
	})

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.StatusText.SetText("Request Saved")
	})

}

func (tui *Tui) handleSendRequest() {
	err := tui.getCurrentRequest()
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ResponseView.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
		})
		return
	}

	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ResponseView.SetText("[yellow]Sending request...[-]")
	})

	resp, err := tui.HttpService.SendRequest(tui.State.CurrentRequest)
	if err != nil {
		tui.Ui.QueueUpdateDraw(func() {
			tui.Components.ResponseView.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
		})
		return
	}
	tui.State.CurrentResponse = resp

	tui.updateOnReceiveResponse()
}

func (tui *Tui) getCurrentRequest() error {
	name := tui.Components.NameInput.GetText()

	_, method := tui.Components.MethodDropdown.GetCurrentOption()

	url := tui.Components.URLInput.GetText()

	headersText := tui.Components.HeadersText.GetText()

	paramsText := tui.Components.ParamsText.GetText()

	_, bodyType := tui.Components.BodyType.GetCurrentOption()

	body := tui.Components.BodyText.GetText()

	newRequest := *domain.NewRequest()

	err := newRequest.BuildRequest(name, method, url, headersText, paramsText, bodyType, body, tui.Config)
	if err != nil {
		return err
	}

	tui.State.CurrentRequest = &newRequest

	return nil
}

func (tui *Tui) loadSavedRequests() {
	if tui.HttpService == nil {
		return
	}
	savedReqs, err := tui.HttpService.GetSavedRequests()
	if err != nil {
		log.Printf("Error loading saved requests: %v", err)
		return
	}

	if len(savedReqs) > 0 {
		tui.State.SavedRequests = savedReqs

		tui.Components.RequestList.Clear()

		for _, req := range tui.State.SavedRequests {
			name := req.Name
			if name == "" {
				name = "unnamed"
			}

			itemText := fmt.Sprintf("%-15s|  %-6s|%s", name, req.Method, req.URL)
			tui.Components.RequestList.AddItem(itemText, "", 0, nil)
		}
	} else {
		tui.Components.RequestList.Clear()
	}
}

func (tui *Tui) updateOnReceiveResponse() {
	responseText := responseStringBuilder(tui.State.CurrentResponse)
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ResponseView.SetText(responseText)
	})
}

func (tui *Tui) populateRequest(req *domain.Request) {
	methodIdx := 0
	bodyTypeIdx := 0

	switch req.Method {
	case "GET":
		methodIdx = 0
	case "POST":
		methodIdx = 1
	case "PUT":
		methodIdx = 2
	case "DELETE":
		methodIdx = 3
	case "HEAD":
		methodIdx = 4
	}

	switch req.ContentType["Content-Type"] {
	case "Text":
		bodyTypeIdx = 0
	case "JSON":
		bodyTypeIdx = 1
	}

	tui.Components.MethodDropdown.SetCurrentOption(methodIdx)
	tui.Components.URLInput.SetText(req.URL)
	tui.Components.NameInput.SetText(req.Name)
	tui.Components.HeadersText.SetText(mapToString(req.Headers), true)
	tui.Components.ParamsText.SetText(mapToString(req.Params), true)
	tui.Components.BodyType.SetCurrentOption(bodyTypeIdx)
	tui.Components.BodyText.SetText(req.Body, true)

}

func responseStringBuilder(resp *domain.Response) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "[yellow]Status:[-] [blue]%s[-]\n\n", resp.Status)
	fmt.Fprintf(&builder, "[yellow]Response time:[-] [blue]%s[-]\n\n", resp.ResponseTime)
	fmt.Fprintf(&builder, "[yellow]Content-Type:[-] [blue]%s[-]\n\n", resp.ContentType)
	fmt.Fprintf(&builder, "[yellow]Content-Length:[-] [blue]%d[-]\n\n", resp.ContentLenght)

	if resp.Body != "" {
		fmt.Fprintf(&builder, "[yellow]Body:[-]\n")
		fmt.Fprint(&builder, resp.Body)
	} else {
		fmt.Fprint(&builder, "[blue]No body[-]")
	}

	return builder.String()
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	var parts []string
	for k, v := range m {
		if k != "User-Agent" {
			parts = append(parts, fmt.Sprintf("%s:%s", k, v))
		}
	}
	return strings.Join(parts, ", ")
}
