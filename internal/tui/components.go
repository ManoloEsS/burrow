package tui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createTuiLayout() *UIComponents {
	components := &UIComponents{}

	topFlex := tview.NewFlex()

	// logo
	components.LogoText = tview.NewTextView().SetText(
		"  _ __                      \n ( /  )                     \n  /--< , , _   _   __ , , , \n /___/(_/_/ (_/ (_(_)(_(_/_ ",
	).SetTextColor(tcell.ColorBlue)

	// keybindings
	components.BindingsText = tview.NewTextView().SetText("C-f: request   | C-s: send req\nC-t: name input| C-e: start server\nC-g: path input| C-x: kill server\nC-l: saved reqs| C-r: reload server").
		SetTextColor(tcell.ColorGray)

	// server status
	serverFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	components.ServerPath = tview.NewInputField()
	components.ServerPath.SetPlaceholder(" path/to/server").
		SetPlaceholderStyle(tcell.StyleDefault).
		SetLabel("server").
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldWidth(26)

	components.ServerStatus = tview.NewTextView()
	components.ServerStatus.SetDynamicColors(true).
		SetText("Server not running")

	serverFlex.AddItem(components.ServerStatus, 0, 1, false).
		AddItem(components.ServerPath, 0, 1, false)

	topFlex.AddItem(components.LogoText, 0, 3, false).
		AddItem(components.BindingsText, 0, 4, false).
		AddItem(serverFlex, 0, 3, false)

	// bottom
	bottomFlex := tview.NewFlex()

	// request builder form
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	components.URLInput = tview.NewInputField()
	components.URLInput.SetPlaceholder("default localhost:8080").
		SetPlaceholderStyle(tcell.StyleDefault).
		SetLabel("URL ").
		SetFieldBackgroundColor(tcell.ColorLightCoral)

	components.HeadersText = tview.NewTextArea()
	components.HeadersText.SetPlaceholder("key:value key:value").
		SetLabel("Headers").
		SetPlaceholderStyle(tcell.StyleDefault).
		SetSize(3, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorBlue)

	components.ParamsText = tview.NewTextArea()
	components.ParamsText.SetPlaceholder("key:value key:value").
		SetLabel("Params").
		SetPlaceholderStyle(tcell.StyleDefault).
		SetSize(3, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorBlue)

	components.BodyText = tview.NewTextArea()
	components.BodyText.SetPlaceholder("Your body content here").
		SetLabel("Body").
		SetPlaceholderStyle(tcell.StyleDefault).
		SetSize(8, 0).
		SetFormAttributes(8, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorBlack, tcell.ColorBlue)

	form := tview.NewForm().
		AddDropDown("Method", []string{"GET", "POST", "PUT", "DELETE", "HEAD"}, 0, nil).
		AddFormItem(components.URLInput).
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

	bodyFormItem := form.GetFormItem(4)
	if bodyDropDown, ok := bodyFormItem.(*tview.DropDown); ok {
		components.BodyType = bodyDropDown
		components.BodyType.SetCurrentOption(0)
	} else {
		log.Printf("Warning: Failed to extract Body type dropdown from form")
	}

	leftFlex.AddItem(form, 0, 1, false)

	// response, list and req input
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	responseFlex := tview.NewFlex()
	components.ResponseView = tview.NewTextView()
	components.ResponseView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Response").
		SetTitleAlign(tview.AlignLeft)

	responseFlex.AddItem(components.ResponseView, 0, 1, false)

	bottomRightFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// request name
	components.NameInput = tview.NewInputField()
	components.NameInput.SetLabel("Request name ")

	// request list section
	components.RequestList = tview.NewList()
	components.RequestList.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle("Saved Requests").
		SetTitleAlign(tview.AlignLeft)

	// add request list and server section to bottom right
	bottomRightFlex.AddItem(components.RequestList, 0, 3, false)
	bottomRightFlex.AddItem(components.NameInput, 0, 2, false)

	rightFlex.AddItem(responseFlex, 0, 8, false).
		AddItem(bottomRightFlex, 0, 2, false)

	bottomFlex.AddItem(leftFlex, 0, 7, false).
		AddItem(rightFlex, 0, 9, false)

	components.MainLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	components.MainLayout.AddItem(topFlex, 0, 2, false).
		AddItem(bottomFlex, 0, 10, false)

	return components
}
