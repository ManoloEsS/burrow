package tui

import (
	"fmt"
	"strings"

	"github.com/ManoloEsS/burrow/internal/domain"
)

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
	_, method := tui.Components.MethodDropdown.GetCurrentOption()

	url := tui.Components.URLInput.GetText()

	headersText := tui.Components.HeadersText.GetText()

	paramsText := tui.Components.ParamsText.GetText()

	_, bodyType := tui.Components.BodyType.GetCurrentOption()

	body := tui.Components.BodyText.GetText()

	newRequest := *domain.NewRequest()

	err := newRequest.BuildRequest(method, url, headersText, paramsText, bodyType, body, tui.Config)
	if err != nil {
		return err
	}

	tui.State.CurrentRequest = &newRequest

	return nil
}

func (tui *Tui) loadSavedRequests() {
	err := tui.HttpService.GetSavedRequests()
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

func (tui *Tui) updateOnReceiveResponse() {
	responseText := tui.responseStringBuilder(tui.State.CurrentResponse)
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.ResponseView.SetText(responseText)
	})
}

func (tui *Tui) updateRequestHistory() {
	err := tui.HttpService.GetSavedRequests()
	if err != nil {
		return
	}

	tui.State.RequestHistory = nil

}

func (tui *Tui) responseStringBuilder(resp *domain.Response) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Status: %s\n\n", resp.Status))
	builder.WriteString(fmt.Sprintf("Response time: %s\n\n", resp.ResponseTime))
	builder.WriteString(fmt.Sprintf("Content-Type: %s\n\n", resp.ContentType))
	builder.WriteString(fmt.Sprintf("Content-Length: %d\n\n", resp.ContentLenght))

	if resp.Body != "" {
		builder.WriteString("Body:\n")
		builder.WriteString(resp.Body)
	} else {
		builder.WriteString(`No body`)
	}

	return builder.String()
}
