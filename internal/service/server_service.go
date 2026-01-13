package service

import (
	"sync"

	"github.com/ManoloEsS/burrow/internal/config"
)

type serverService struct {
	config        *config.Config
	currentStatus ServerStatus
	statusMutex   sync.RWMutex
	serverCancel  chan struct{}
	serverWG      sync.WaitGroup
	crashDetected bool
}

func NewServerService() ServerService {
	return &serverService{
		currentStatus: ServerStatus{
			Running: false,
			Status:  "Server not running",
		},
		serverCancel: make(chan struct{}),
	}
}

func (s *serverService) StartServer(path string) error {
	return nil
}

func (s *serverService) StopServer() error {
	return nil
}

func (s *serverService) stopServerInternal() error {
	return nil
}

func (s *serverService) GetStatus() ServerStatus {
	return s.currentStatus
}

func (s *serverService) runServer(path string) {
}

func (s *serverService) handleServerCrash() {
}
