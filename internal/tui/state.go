package tui

import (
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

type UIState struct {
	CurrentRequest        *domain.Request
	CurrentServer         service.ServerStatus
	RequestHistory        []*domain.Request
	CurrentResponse       *domain.Response
	CurrentFormFocusIndex int
	CurrentFocused        tview.Primitive
}

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
}
