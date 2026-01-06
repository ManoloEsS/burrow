package service

import (
	"fmt"
	"sync"

	"github.com/ManoloEsS/burrow/internal/config"
)

type serverService struct {
	config         *config.Config
	updateCallback ServerUpdateCallback
	currentStatus  ServerStatus
	statusMutex    sync.RWMutex
	serverCancel   chan struct{}
	serverWG       sync.WaitGroup
	crashDetected  bool
}

func NewServerService(config *config.Config, callback ServerUpdateCallback) ServerService {
	return &serverService{
		currentStatus: ServerStatus{
			Running: false,
			Status:  "Server not running",
		},
		serverCancel:   make(chan struct{}),
		config:         config,
		updateCallback: callback,
	}
}

func (s *serverService) StartServer(path string) error {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	if s.currentStatus.Running {
		return fmt.Errorf("server is already running")
	}

	s.stopServerInternal()

	s.serverCancel = make(chan struct{})

	s.currentStatus = ServerStatus{
		Running: true,
		Path:    path,
		Status:  fmt.Sprintf("Server running on %s", path),
	}
	s.crashDetected = false

	s.serverWG.Add(1)
	go s.runServer(path)

	if s.updateCallback != nil {
		s.updateCallback(s.currentStatus)
	}

	return nil
}

func (s *serverService) StopServer() error {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	return s.stopServerInternal()
}

func (s *serverService) stopServerInternal() error {
	if !s.currentStatus.Running {
		return fmt.Errorf("server is not running")
	}

	close(s.serverCancel)

	s.serverWG.Wait()

	s.currentStatus = ServerStatus{
		Running: false,
		Status:  "Server stopped",
	}

	if s.updateCallback != nil {
		s.updateCallback(s.currentStatus)
	}

	return nil
}

func (s *serverService) GetStatus() ServerStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()
	return s.currentStatus
}

func (s *serverService) runServer(path string) {
	defer s.serverWG.Done()

	for {
		select {
		case <-s.serverCancel:
			// Server stopped gracefully
			return

		}
	}
}

func (s *serverService) handleServerCrash() {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	s.crashDetected = true
	s.currentStatus = ServerStatus{
		Running: false,
		Status:  "Server crashed",
	}

	if s.updateCallback != nil {
		s.updateCallback(s.currentStatus)
	}
}
