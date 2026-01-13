package service

import (
	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
)

type Services struct {
	HttpClientService HttpClientService
	ServerService     ServerService
	Config            *config.Config
	Database          *database.Database
}

func NewServices(database *database.Database, config *config.Config) *Services {
	return &Services{
		Config:            config,
		Database:          database,
		HttpClientService: NewHttpClientService(database, config),
		ServerService:     NewServerService(config),
	}
}
