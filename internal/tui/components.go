package tui

import (
	"fmt"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIComponents struct {
	MainLayout       *tview.Flex
	Pages            *tview.Pages
	Form             *tview.Form
	LogoText         *tview.TextView
	BindingsText     *tview.TextView
	KeybindingsModal *tview.Flex
	ServerStatus     *tview.TextView
	ServerPath       *tview.InputField

	MethodDropdown *tview.DropDown
	URLInput       *tview.InputField
	HeadersText    *tview.TextArea
	ParamsText     *tview.TextArea
	BodyText       *tview.TextArea
	BodyType       *tview.DropDown

	ResponseView *tview.TextView

	RequestList *tview.List
	NameInput   *tview.InputField
	StatusText  *tview.TextView
}

func createTuiLayout(cfg *config.Config) *UIComponents {
	components := &UIComponents{}

	components.createLogoComponent()

	components.createKeybindingsComponent()

	components.createKeybindingsModal()

	components.createServerPathComponent()

	components.createServerStatusComponent()

	components.createUrlInputComponent(cfg)

	components.createHeadersTextComponent()

	components.createParamsTextComponent()

	components.createBodyTextComponent()

	components.createResponseViewComponent()

	components.createNameInputComponent()

	components.createRequestListComponent()

	components.createFormAndSetup()

	components.createStatusComponent()

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	leftFlex.AddItem(components.LogoText, 5, 1, false).
		AddItem(components.Form, 0, 1, true)

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(components.ServerStatus, 1, 1, false).
		AddItem(components.ServerPath, 1, 1, false).
		AddItem(components.StatusText, 1, 1, false).
		AddItem(components.ResponseView, 0, 3, false).
		AddItem(components.RequestList, 0, 1, false)

	components.MainLayout = tview.NewFlex().SetDirection(tview.FlexColumn)
	components.MainLayout.AddItem(leftFlex, 0, 7, true).
		AddItem(rightFlex, 0, 9, false)

	components.Pages = tview.NewPages()
	components.Pages.AddPage("main", components.MainLayout, true, true)
	components.Pages.AddPage("keybindings", components.KeybindingsModal, true, false)

	return components
}

func (components *UIComponents) createUrlInputComponent(cfg *config.Config) {
	components.URLInput = tview.NewInputField()
	components.URLInput.SetPlaceholder(fmt.Sprintf("default localhost:%s", cfg.App.DefaultPort)).
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey)).
		SetPlaceholderTextColor(tcell.ColorBlue).
		SetLabel("URL ").
		SetFieldBackgroundColor(tcell.ColorLightCoral)
}

func (components *UIComponents) createNameInputComponent() {
	components.NameInput = tview.NewInputField()
	components.NameInput.SetPlaceholder("name to be saved as").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey)).
		SetPlaceholderTextColor(tcell.ColorBlue).
		SetLabel("Name ").
		SetFieldBackgroundColor(tcell.ColorLightCoral)
}

func (components *UIComponents) createFormAndSetup() {
	form := tview.NewForm().
		AddDropDown("Method", []string{"GET", "POST", "PUT", "DELETE", "HEAD"}, 0, nil).
		AddFormItem(components.URLInput).
		AddFormItem(components.NameInput).
		AddFormItem(components.HeadersText).
		AddFormItem(components.ParamsText).
		AddDropDown("Body", []string{"Text", "JSON"}, 0, nil).
		AddFormItem(components.BodyText)

	form.SetFieldTextColor(tcell.ColorBlack)
	form.ClearButtons().SetButtonTextColor(tcell.ColorBlack).
		SetItemPadding(1)

	components.Form = form

	methodFormItem := form.GetFormItem(0)
	if methodDropDown, ok := methodFormItem.(*tview.DropDown); ok {
		components.MethodDropdown = methodDropDown
		components.MethodDropdown.SetCurrentOption(0)
	}

	bodyFormItem := form.GetFormItem(5)
	if bodyDropDown, ok := bodyFormItem.(*tview.DropDown); ok {
		components.BodyType = bodyDropDown
		components.BodyType.SetCurrentOption(0)
	}

}

func (components *UIComponents) createLogoComponent() {
	components.LogoText = tview.NewTextView().SetText(
		"┳┓           \n┣┫┓┏┏┓┏┓┏┓┓┏┏\n┻┛┗┻┛ ┛ ┗┛┗┻┛",
	).SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createStatusComponent() {
	components.StatusText = tview.NewTextView().SetText("Ready!").
		SetWrap(true).
		SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createKeybindingsComponent() {
	components.BindingsText = tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[white]Request form[-]     [blue]|[-][-][white]Response view[-]        [blue]|[-][white]Saved requests list[-][blue]|[-][white]Server[-]
C-f: focus form  [blue]|[-] C-t: focus resp     [blue]|[-] C-l: focus list   [blue]|[-] C-g: focus input
C-s: send request[blue]|[-] j/k:scroll    ↑↓    [blue]|[-] j/k:navigate  ↑↓  [blue]|[-] C-x: kill server
C-a: save request[blue]|[-][blue]_____________________|[-] C-o: load request [blue]|[-] C-r: start server
C-n/p: navigate↑↓  C-u: clear form     [blue]|[-] C-d: del request  [blue]|[-]`).
		SetTextColor(tcell.ColorGray)
}

func (components *UIComponents) createKeybindingsModal() {
	modalText := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[white]Left side[-]                [blue]|[-][white]Right side[-]
Logo + Form              [blue]|[-] Server + Response + Requests
M-h: left    M-l: right  [blue]|[-] M-h: left   M-l: right
M-j: down   M-k: up      [blue]|[-] M-j: down  M-k: up

[white]Actions[-]
C-s: send request  C-a: save request  C-r: start server
C-x: kill server  C-u: clear form     C-d: del request
C-o: load request

[white]Inside components[-]
j/k: navigate list  C-n/p: navigate form  j/k: scroll response

Press [yellow]Esc[-] or [yellow]Alt+I[-] to close`).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignCenter)

	components.KeybindingsModal = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(modalText, 80, 1, true).
			AddItem(nil, 0, 1, false), 12, 1, false).
		AddItem(nil, 0, 1, false)

	components.KeybindingsModal.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" Keybindings ").
		SetTitleColor(tcell.ColorYellow).
		SetBackgroundColor(tcell.ColorBlack)
}

func (components *UIComponents) createServerPathComponent() {
	components.ServerPath = tview.NewInputField()
	components.ServerPath.SetPlaceholder("path/to/server").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey)).
		SetPlaceholderTextColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorBlack)
}

func (components *UIComponents) createServerStatusComponent() {
	components.ServerStatus = tview.NewTextView()
	components.ServerStatus.SetDynamicColors(true).
		SetWrap(true).
		SetText("Server not running")
}

func (components *UIComponents) createHeadersTextComponent() {
	components.HeadersText = tview.NewTextArea()
	components.HeadersText.SetPlaceholder("key:value, key:value").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorBlue)).
		SetLabel("Headers").
		SetSize(2, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorLightCoral)
}

func (components *UIComponents) createParamsTextComponent() {
	components.ParamsText = tview.NewTextArea()
	components.ParamsText.SetPlaceholder("key:value, key:value").
		SetLabel("Params").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorBlue)).
		SetSize(2, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorLightCoral)
}

func (components *UIComponents) createBodyTextComponent() {
	components.BodyText = tview.NewTextArea()
	components.BodyText.SetPlaceholder("Your body content here").
		SetLabel("Body").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorBlue)).
		SetSize(10, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorLightCoral)
}

func (components *UIComponents) createResponseViewComponent() {
	components.ResponseView = tview.NewTextView()
	components.ResponseView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Response").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorYellow)
}

func (components *UIComponents) createRequestListComponent() {
	components.RequestList = tview.NewList()
	components.RequestList.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle("Saved Requests").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorYellow)
}
