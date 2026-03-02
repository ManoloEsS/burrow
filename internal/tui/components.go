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
	FormTitle        *tview.TextView
	LogoText         *tview.TextView
	BindingsText     *tview.TextView
	KeybindingsModal *tview.Flex
	InfoText         *tview.TextView // Added: accessible for updates
	ServerStatus     *tview.TextView
	ServerPath       *tview.InputField
	ServerInfoBox    *tview.Flex
	ResponseView     *tview.TextView
	RequestList      *tview.List

	MethodDropdown *tview.DropDown
	URLInput       *tview.InputField
	HeadersText    *tview.TextArea
	ParamsText     *tview.TextArea
	BodyText       *tview.TextArea
	BodyType       *tview.DropDown

	NameInput  *tview.InputField
	StatusText *tview.TextView
}

func createTuiLayout(cfg *config.Config) *UIComponents {
	components := &UIComponents{}

	components.createLogoComponent()

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

	serverInfoBox := tview.NewFlex().SetDirection(tview.FlexRow)
	serverInfoBox.AddItem(components.ServerStatus, 1, 1, false).
		AddItem(components.ServerPath, 1, 1, false).
		AddItem(components.StatusText, 1, 1, false)
	serverInfoBox.SetBorder(true).
		SetTitle(" Status ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorBlue)
	components.ServerInfoBox = serverInfoBox

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	leftFlex.AddItem(components.LogoText, 3, 1, false).
		// AddItem(components.FormTitle, 1, 1, false).
		AddItem(components.Form, 0, 1, true)

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(components.ServerInfoBox, 5, 1, false).
		AddItem(components.ResponseView, 0, 4, false).
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

	components.FormTitle = tview.NewTextView().
		SetText(" Request ").
		SetTextColor(tcell.ColorBlue)

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
		" ┳┓           \n ┣┫┓┏┏┓┏┓┏┓┓┏┏\n ┻┛┗┻┛ ┛ ┗┛┗┻┛",
	).SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createStatusComponent() {
	components.StatusText = tview.NewTextView().SetText("Ready!").
		SetWrap(true).
		SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createKeybindingsModal() {
	infoView := tview.NewTextView()
	infoView.SetDynamicColors(true)
	infoView.SetText(`hello`)
	infoView.SetBorder(true)
	infoView.SetBorderColor(tcell.ColorYellow)
	infoView.SetTitle(" Keybindings ")
	infoView.SetTitleColor(tcell.ColorBlue)

	components.InfoText = infoView // Added: store reference for updates

	components.KeybindingsModal = tview.NewFlex(). // Changed: FlexColumn for vertical stacking
							AddItem(nil, 0, 1, false).
							AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
								AddItem(nil, 0, 1, false).
								AddItem(infoView, 20, 1, true).
								AddItem(nil, 0, 1, false), 60, 1, false).
							AddItem(nil, 0, 1, false)

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
		SetTitle(" Response ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorBlue)
}

func (components *UIComponents) createRequestListComponent() {
	components.RequestList = tview.NewList()
	components.RequestList.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(" Saved Requests ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorBlue)
}
