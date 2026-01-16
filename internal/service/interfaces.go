package service

import (
	"github.com/ManoloEsS/burrow/internal/domain"
)

type ServerStatus struct {
	Running bool
	Path    string
	Status  string
}

type HttpClientService interface {
	SendRequest(*domain.Request) (*domain.Response, error)
	SaveRequest(*domain.Request) error
	DeleteRequest(string) error
	GetSavedRequests() ([]*domain.Request, error)
}

type ServerService interface {
	StartServer(string) error
	StopServer() error
	HealthCheck() ServerStatus
}
