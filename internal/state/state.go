package state

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/domain"
)

type State struct {
	Cfg      *config.Config
	Screen   ScreenState
	Requests domain.Request
}

type ScreenState int

const (
	MainScreen ScreenState = iota
	NewReqScreen
	RetrieveScreen
)
