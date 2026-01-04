package service

import (
	"fmt"
	"sync"

	"github.com/ManoloEsS/burrow/internal/config"
)

type serverService struct {
	config         *config.Config
	updateCallback func(ServerStatus)
	currentStatus  ServerStatus
	statusMutex    sync.RWMutex
	serverCancel   chan struct{}
	serverWG       sync.WaitGroup
	crashDetected  bool
}

// Starts module in service layer to handle server behaviors
func NewServerService(config *config.Config) ServerService {
	return &serverService{
		currentStatus: ServerStatus{
			Running: false,
			Status:  "Server not running",
		},
		serverCancel: make(chan struct{}),
		config:       config,
	}
}

// starts server go routine and updates ui
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

// Stops server instance
func (s *serverService) StopServer() error {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	return s.stopServerInternal()
}

func (s *serverService) stopServerInternal() error {
	if !s.currentStatus.Running {
		return fmt.Errorf("server is not running")
	}

	// Cancel server goroutine
	close(s.serverCancel)

	// Wait for server to stop
	s.serverWG.Wait()

	// Update status
	s.currentStatus = ServerStatus{
		Running: false,
		Status:  "Server stopped",
	}

	// Notify UI
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

func (s *serverService) SetUpdateCallback(callback func(ServerStatus)) {
	s.updateCallback = callback
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

	// Notify UI of crash
	if s.updateCallback != nil {
		s.updateCallback(s.currentStatus)
	}
}
