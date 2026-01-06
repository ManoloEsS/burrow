package service

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
)

type Services struct {
	RequestService RequestService
	ServerService  ServerService
	Config         *config.Config
	Database       *database.Database
}

func NewServices(database *database.Database, config *config.Config,
	requestCallback RequestUpdateCallBack,
	serverCallback ServerUpdateCallback) *Services {
	return &Services{
		Config:         config,
		Database:       database,
		RequestService: NewRequestService(database, config, requestCallback),
		ServerService:  NewServerService(config, serverCallback),
	}
}
