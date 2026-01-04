package service

import (
	"time"

	"github.com/ManoloEsS/burrow/internal/domain"
)

type Response struct {
	ID         int
	RequestID  int
	StatusCode int
	Body       string
	Created    time.Time
}

// Tracks status of server for updating ui
type ServerStatus struct {
	Running bool
	Path    string
	Status  string
}

// Service module for handling http request behaviors
type RequestService interface {
	SendRequest(req *domain.Request) (*Response, error)
	SaveRequest(req *domain.Request) error
	GetSavedRequests() error
	SetUpdateCallback(callback func(*Response))
}

// Service module for handling server behaviors
type ServerService interface {
	StartServer(path string) error
	StopServer() error
	GetStatus() ServerStatus
	SetUpdateCallback(callback func(ServerStatus))
}
