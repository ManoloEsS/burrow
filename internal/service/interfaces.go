package service

import (
	"github.com/ManoloEsS/burrow/internal/domain"
)

type HttpClientService interface {
	SendRequest(*domain.Request) (*domain.Response, error)
	SaveRequest(*domain.Request) error
	DeleteRequest(string) error
	GetSavedRequests() ([]*domain.Request, error)
}

type ServerService interface {
	StartServer(path string, port string, updateChan chan UIEvent) error
	StopServer() error
}
