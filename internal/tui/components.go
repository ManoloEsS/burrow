package tui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIComponents struct {
	MainLayout   *tview.Flex
	Form         *tview.Form
	LogoText     *tview.TextView
	BindingsText *tview.TextView
	ServerStatus *tview.TextView
	ServerPath   *tview.InputField

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

func createTuiLayout() *UIComponents {
	components := &UIComponents{}

	components.createLogoComponent()

	components.createKeybindingsComponent()

	components.createServerPathComponent()

	components.createServerStatusComponent()

	components.createUrlInputComponent()

	components.createHeadersTextComponent()

	components.createParamsTextComponent()

	components.createBodyTextComponent()

	components.createResponseViewComponent()

	components.createNameInputComponent()

	components.createRequestListComponent()

	components.createFormAndSetup()

	components.createStatusComponent()

	topFlex := tview.NewFlex()

	serverFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	serverFlex.AddItem(components.ServerStatus, 0, 1, false).
		AddItem(components.ServerPath, 0, 1, false).
		AddItem(components.StatusText, 0, 1, false)

	topFlex.AddItem(components.LogoText, 0, 3, false).
		AddItem(components.BindingsText, 0, 4, false)

	bottomFlex := tview.NewFlex()

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	leftFlex.AddItem(components.Form, 0, 1, false)

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	responseFlex := tview.NewFlex()

	responseFlex.AddItem(components.ResponseView, 0, 1, false)

	bottomRightFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	bottomRightFlex.AddItem(components.RequestList, 0, 3, false)
	bottomRightFlex.AddItem(serverFlex, 0, 2, false)

	rightFlex.AddItem(responseFlex, 0, 8, false).
		AddItem(bottomRightFlex, 0, 2, false)

	bottomFlex.AddItem(leftFlex, 0, 7, false).
		AddItem(rightFlex, 0, 9, false)

	components.MainLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	components.MainLayout.AddItem(topFlex, 0, 2, false).
		AddItem(bottomFlex, 0, 10, false)

	return components
}

func (components *UIComponents) createUrlInputComponent() {
	components.URLInput = tview.NewInputField()
	components.URLInput.SetPlaceholder("default localhost:8080").
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
	} else {
		log.Printf("Warning: Failed to extract Body type dropdown from form")
	}

}

func (components *UIComponents) createLogoComponent() {
	components.LogoText = tview.NewTextView().SetText(
		"  _ __                      \n ( /  )                     \n  /--< , , _   _   __ , , , \n /___/(_/_/ (_/ (_(_)(_(_/_ ",
	).SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createStatusComponent() {
	components.StatusText = tview.NewTextView().SetText("Ready!").
		SetTextColor(tcell.ColorBlue)
}

func (components *UIComponents) createKeybindingsComponent() {
	components.BindingsText = tview.NewTextView().SetText("C-f: request   | C-s: send req\nC-t: name input| C-e: start server\nC-g: path input| C-x: kill server\nC-l: saved reqs| C-r: reload server\n\nDropdowns: j/k navigate, Enter to select").
		SetTextColor(tcell.ColorGray)
}

func (components *UIComponents) createServerPathComponent() {
	components.ServerPath = tview.NewInputField()
	components.ServerPath.SetPlaceholder(" path/to/server").
		SetLabel("server").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey)).
		SetPlaceholderTextColor(tcell.ColorBlue).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldWidth(26)
}

func (components *UIComponents) createServerStatusComponent() {
	components.ServerStatus = tview.NewTextView()
	components.ServerStatus.SetDynamicColors(true).
		SetText("Server not running")
}

func (components *UIComponents) createHeadersTextComponent() {
	components.HeadersText = tview.NewTextArea()
	components.HeadersText.SetPlaceholder("key:value key:value").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorBlue)).
		SetLabel("Headers").
		SetSize(2, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorLightCoral)
}

func (components *UIComponents) createParamsTextComponent() {
	components.ParamsText = tview.NewTextArea()
	components.ParamsText.SetPlaceholder("key:value key:value").
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
		SetSize(8, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorLightCoral)
}

func (components *UIComponents) createResponseViewComponent() {
	components.ResponseView = tview.NewTextView()
	components.ResponseView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Response").
		SetTitleAlign(tview.AlignLeft)
}

func (components *UIComponents) createRequestListComponent() {
	components.RequestList = tview.NewList()
	components.RequestList.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle("Saved Requests").
		SetTitleAlign(tview.AlignLeft)
}
