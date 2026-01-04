package service

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
)

// Service layer with request and server modules
type Services struct {
	RequestService RequestService
	ServerService  ServerService
	Config         *config.Config
}

// initialize service layer and modules
func NewServices(database *database.Database, config *config.Config) *Services {
	return &Services{
		RequestService: NewRequestService(database, config),
		ServerService:  NewServerService(config),
	}
}
