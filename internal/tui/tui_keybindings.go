package tui

import "github.com/gdamore/tcell/v2"

func (tui *Tui) setupKeybindings() {
	tui.Ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// handlers
		case tcell.KeyCtrlS:
			go tui.handleSendRequest()
			return nil
		case tcell.KeyCtrlQ:
			tui.Stop()
			return nil
		case tcell.KeyCtrlA:
			go tui.handleSaveRequest()
			return nil
		case tcell.KeyCtrlR:
			go tui.handleStartServer()
			return nil
		case tcell.KeyCtrlX:
			go tui.handleStopServer()
			return nil
		case tcell.KeyCtrlD:
			if tui.State.CurrentFocused == tui.Components.RequestList {
				go tui.handleDeleteRequest()
				return nil
			}
			return nil
		case tcell.KeyCtrlO:
			if tui.State.CurrentFocused == tui.Components.RequestList {
				go tui.handleLoadRequest()
				return nil
			}
			return nil
		case tcell.KeyCtrlU:
			if tui.State.CurrentFocused == tui.Components.Form {
				go tui.clear()
				return nil
			}
			return nil

		// navigation
		case tcell.KeyEsc:
			if tui.State.KeybindingsVisible {
				tui.toggleKeybindingsModal()
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h':
				if event.Modifiers()&tcell.ModAlt != 0 {
					tui.focusLeft()
					return nil
				}
				return event
			case 'l':
				if event.Modifiers()&tcell.ModAlt != 0 {
					tui.focusRight()
					return nil
				}
				return event
			case 'j':
				if event.Modifiers()&tcell.ModAlt != 0 {
					tui.focusDown()
					return nil
				}
				if tui.State.CurrentFocused == tui.Components.RequestList {
					tui.navigateList(1)
					return nil
				}
				return event
			case 'k':
				if event.Modifiers()&tcell.ModAlt != 0 {
					tui.focusUp()
					return nil
				}
				if tui.State.CurrentFocused == tui.Components.RequestList {
					tui.navigateList(-1)
					return nil
				}
				return event
			case 'i':
				if event.Modifiers()&tcell.ModAlt != 0 {
					tui.toggleKeybindingsModal()
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

func (tui *Tui) focusForm() {
	tui.State.CurrentSide = "left"
	tui.State.CurrentFocused = tui.Components.Form
	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
	tui.updateBorderColors()
}

func (tui *Tui) navigateForm(forward bool) {
	subcompCount := tui.Components.Form.GetFormItemCount()
	if forward {
		tui.State.CurrentFormFocusIndex = (tui.State.CurrentFormFocusIndex + 1) % subcompCount
	} else {
		tui.State.CurrentFormFocusIndex = (tui.State.CurrentFormFocusIndex - 1 + 7) % subcompCount
	}

	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
}

func (tui *Tui) focusSpecificFormComponent(index int) {
	component := tui.Components.Form.GetFormItem(index)
	if component != nil {
		tui.Ui.SetFocus(component)
	}
}

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

func (tui *Tui) clear() {
	tui.Ui.QueueUpdateDraw(func() {
		tui.Components.MethodDropdown.SetCurrentOption(0)
		tui.Components.URLInput.SetText("")
		tui.Components.NameInput.SetText("")
		tui.Components.HeadersText.SetText("", true)
		tui.Components.ParamsText.SetText("", true)
		tui.Components.BodyType.SetCurrentOption(0)
		tui.Components.BodyText.SetText("", true)
	})
}

func (tui *Tui) toggleKeybindingsModal() {
	tui.State.KeybindingsVisible = !tui.State.KeybindingsVisible
	if tui.State.KeybindingsVisible {
		tui.State.LastFocusedBeforeModal = tui.State.CurrentFocused
		tui.Components.Pages.ShowPage("keybindings")
	} else {
		tui.Components.Pages.HidePage("keybindings")
		if tui.State.LastFocusedBeforeModal != nil {
			tui.Ui.SetFocus(tui.State.LastFocusedBeforeModal)
		}
	}
}

func (tui *Tui) focusLeft() {
	tui.State.CurrentSide = "left"
	tui.State.CurrentFocused = tui.Components.Form
	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
	tui.updateBorderColors()
}

func (tui *Tui) focusRight() {
	tui.State.CurrentSide = "right"
	tui.focusRightComponent(tui.State.CurrentRightComponentIndex)
	tui.updateBorderColors()
}

func (tui *Tui) focusDown() {
	if tui.State.CurrentSide == "left" {
		tui.navigateForm(true)
		return
	}

	tui.State.CurrentRightComponentIndex = (tui.State.CurrentRightComponentIndex + 1) % 3
	tui.focusRightComponent(tui.State.CurrentRightComponentIndex)
}

func (tui *Tui) focusUp() {
	if tui.State.CurrentSide == "left" {
		tui.navigateForm(false)
		return
	}

	tui.State.CurrentRightComponentIndex = (tui.State.CurrentRightComponentIndex - 1 + 3) % 3
	tui.focusRightComponent(tui.State.CurrentRightComponentIndex)
}

func (tui *Tui) focusRightComponent(index int) {
	switch index {
	case 0:
		tui.State.CurrentFocused = tui.Components.ServerPath
		tui.Ui.SetFocus(tui.Components.ServerPath)
	case 1:
		tui.State.CurrentFocused = tui.Components.ResponseView
		tui.Ui.SetFocus(tui.Components.ResponseView)
	case 2:
		tui.State.CurrentFocused = tui.Components.RequestList
		tui.Ui.SetFocus(tui.Components.RequestList)
	}
	tui.updateBorderColors()
}

func (tui *Tui) updateBorderColors() {
	blue := tcell.ColorBlue
	yellow := tcell.ColorYellow

	tui.Components.FormTitle.SetTextColor(blue)
	tui.Components.ServerInfoBox.SetTitleColor(blue)
	tui.Components.ResponseView.SetTitleColor(blue)
	tui.Components.RequestList.SetTitleColor(blue)

	switch tui.State.CurrentFocused {
	case tui.Components.Form:
		tui.Components.FormTitle.SetTextColor(yellow)
	case tui.Components.ServerPath, tui.Components.ServerStatus, tui.Components.StatusText:
		tui.Components.ServerInfoBox.SetTitleColor(yellow)
	case tui.Components.ResponseView:
		tui.Components.ResponseView.SetTitleColor(yellow)
	case tui.Components.RequestList:
		tui.Components.RequestList.SetTitleColor(yellow)
	}
}
