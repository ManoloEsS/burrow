package state

type ScreenState struct {
	CurrentScreen string
	SelectedID    string
}

func (s *ScreenState) SetScreen(screenID string) {
	knownScreens := map[string]bool{
		"main":     true,
		"request":  true,
		"retrieve": true,
	}

	if knownScreens[screenID] {
		s.CurrentScreen = screenID
	} else {
		s.CurrentScreen = "main"
	}
}
