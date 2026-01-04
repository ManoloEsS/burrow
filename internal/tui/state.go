package tui

import (
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/rivo/tview"
)

// UIState acts as a request, server and saved requests cache
type UIState struct {
	CurrentRequest  *domain.Request
	CurrentServer   service.ServerStatus
	RequestHistory  []*domain.Request
	CurrentResponse *service.Response // In-memory only
}

// UIComponents holds the components that make up the tui
type UIComponents struct {
	MainLayout *tview.Flex
	Form       *tview.Form
	// Top section
	LogoText     *tview.TextView
	BindingsText *tview.TextView
	ServerStatus *tview.TextView
	ServerPath   *tview.InputField

	// Request form section
	MethodDropdown *tview.DropDown
	URLInput       *tview.InputField
	HeadersText    *tview.TextArea
	ParamsText     *tview.TextArea
	BodyText       *tview.TextArea
	BodyType       *tview.DropDown

	// Response section
	ResponseView *tview.TextView

	// Request list section
	RequestList *tview.List
	NameInput   *tview.InputField
}
