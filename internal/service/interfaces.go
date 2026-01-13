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
	SendRequest(req *domain.Request) (*domain.Response, error)
	SaveRequest(req *domain.Request) error
	GetSavedRequests() error
}

type ServerService interface {
	StartServer(path string) error
	StopServer() error
	GetStatus() ServerStatus
}
