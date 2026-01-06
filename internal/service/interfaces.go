package service

import (
	"github.com/ManoloEsS/burrow/internal/domain"
)

type RequestUpdateCallBack func(*domain.Response)
type ServerUpdateCallback func(ServerStatus)

type ServerStatus struct {
	Running bool
	Path    string
	Status  string
}

type RequestService interface {
	SendRequest(req *domain.Request) (*domain.Response, error)
	SaveRequest(req *domain.Request) error
	GetSavedRequests() error
}

type ServerService interface {
	StartServer(path string) error
	StopServer() error
	GetStatus() ServerStatus
}
