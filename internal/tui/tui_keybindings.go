package tui

import "github.com/gdamore/tcell/v2"

func (tui *Tui) setupKeybindings() {
	tui.Ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// events
		case tcell.KeyCtrlS:
			go tui.handleSendRequest()
			return nil
		case tcell.KeyCtrlQ:
			tui.Ui.Stop()
			return nil
		case tcell.KeyCtrlA:
			go tui.handleSaveRequest()
			return nil
		case tcell.KeyF5:
			go tui.handleStartServer()
			return nil
		case tcell.KeyF6:
			go tui.handleStopServer()
			return nil

		// navigation
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
		case tcell.KeyCtrlU:
			if tui.State.CurrentFocused == tui.Components.Form {
				go tui.clear()
				return nil
			}
			return nil
		case tcell.KeyCtrlP:
			if tui.State.CurrentFocused == tui.Components.Form {
				tui.navigateForm(false)
			}
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

func (tui *Tui) focusForm() {
	tui.State.CurrentFocused = tui.Components.Form
	tui.focusSpecificFormComponent(tui.State.CurrentFormFocusIndex)
}

func (tui *Tui) focusServerInput() {
	tui.State.CurrentFocused = tui.Components.ServerPath
	tui.Ui.SetFocus(tui.Components.ServerPath)
}

func (tui *Tui) focusRequestNameInput() {
	tui.State.CurrentFocused = tui.Components.NameInput
	tui.Ui.SetFocus(tui.Components.NameInput)
}

func (tui *Tui) focusRequestList() {
	tui.State.CurrentFocused = tui.Components.RequestList
	tui.Ui.SetFocus(tui.Components.RequestList)
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
