package state

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
)

type State struct {
	Cfg      *config.Config
	DB       *database.Database
	Screen   ScreenState
	Requests map[string]domain.Request
}
