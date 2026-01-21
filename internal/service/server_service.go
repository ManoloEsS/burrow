package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type serverService struct {
	isRunning      bool
	updateChan     chan UIEvent
	serverMu       sync.Mutex
	cancelFunc     context.CancelFunc
	serverProcess  *exec.Cmd
	pathToServer   string
	healthCheckURL string
	httpClient     *http.Client
}

type UIEvent struct {
	Type    string
	Message string
}

func NewServerService() ServerService {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &serverService{
		httpClient: client,
	}
}

func (s *serverService) StartServer(path string, port string, updateChan chan UIEvent) error {
	s.serverMu.Lock()
	defer s.serverMu.Unlock()
	// FIX: Single assignment of updateChan instead of duplicate assignments
	s.updateChan = updateChan
	s.sendEvent("update", "starting server...")
	validPath, err := s.validatePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}

	if s.isRunning {
		return fmt.Errorf("server already running")
	}

	s.pathToServer = validPath
	s.healthCheckURL = "http://localhost:" + port + "/health"
	// FIX: Removed duplicate assignment of updateChan (was set on line 44)

	orchestratorCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	go s.orchestrator(orchestratorCtx)
	s.sendEvent("update", "orchestrator starting...")
	return nil
}

func (s *serverService) orchestrator(ctx context.Context) {
	healthCheckerCtx, healthCheckerCancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	cmd, err := s.runCmdFromPath(healthCheckerCtx)
	if err != nil {
		s.sendEvent("error", fmt.Sprintf("couldn't run file: %v", err))
		defer healthCheckerCancel()
		return
	}
	s.serverMu.Lock()
	s.serverProcess = cmd
	s.isRunning = true
	s.serverMu.Unlock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.healthChecker(healthCheckerCtx)
	}()

	for {
		select {
		case <-ctx.Done():
			// FIX: Improved shutdown sequence for better context cancellation handling
			// Cancel health checker first to prevent interference with server shutdown
			healthCheckerCancel()
			// Wait for health checker goroutine to finish before shutting down server
			wg.Wait()
			// Now safely shutdown the server process
			s.gracefulShutdown()

			s.serverMu.Lock()
			defer s.serverMu.Unlock()

			s.isRunning = false
			s.serverProcess = nil
			return
		}
	}

}

func (s *serverService) runCmdFromPath(ctx context.Context) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "go", "run", s.pathToServer)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *serverService) healthChecker(ctx context.Context) {
	time.Sleep(time.Second * 1)

	s.sendEvent("update", "trying to reach server")

	resp, err := s.httpClient.Get(s.healthCheckURL)
	// FIX: Close response body immediately to prevent resource leak
	// Previous defer would keep all response bodies open until function exit
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil || resp.StatusCode != 200 {
		if err != nil {
			s.sendEvent("error", fmt.Sprintf("cant reach server: %v", err))
		} else {
			s.sendEvent("error", fmt.Sprintf("server returned status %d (expected 200)", resp.StatusCode))
		}
	} else {
		s.sendEvent("update", "server reached successfully")
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := s.httpClient.Get(s.healthCheckURL)
			// FIX: Close response body immediately to prevent resource leak in ticker loop
			// Previous defer would accumulate open file descriptors every 5 seconds
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			if err != nil || resp.StatusCode != 200 {
				if err != nil {
					s.sendEvent("error", fmt.Sprintf("cant reach server: %v", err))
				} else {
					s.sendEvent("error", fmt.Sprintf("server returned status %d (expected 200)", resp.StatusCode))
				}
			}

			s.sendEvent("update", "server healthy")
		}
	}

}

func (s *serverService) StopServer() error {
	s.serverMu.Lock()
	defer s.serverMu.Unlock()

	if s.cancelFunc == nil || !s.isRunning {
		return errors.New("server not running")
	}

	cancel := s.cancelFunc
	s.cancelFunc = nil
	cancel()
	return nil
}

func (s *serverService) gracefulShutdown() {
	// FIX: Single nil check to prevent duplicate validation
	// Previous code checked the same condition twice and sent duplicate messages
	if s.serverProcess == nil || s.serverProcess.Process == nil {
		s.sendEvent("update", "no server process to stop")
		return
	}

	// FIX: Single shutdown message instead of duplicate
	s.sendEvent("update", "stopping server service - context cancelled, waiting for graceful shutdown")

	// Store process reference for shutdown operations
	process := s.serverProcess

	done := make(chan error, 1)
	go func() {
		done <- process.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			s.sendEvent("error", fmt.Sprintf("server process exited with error: %v", err))
		} else {
			s.sendEvent("update", "server process shut down gracefully")
		}
	case <-time.After(5 * time.Second):
		s.sendEvent("error", "server process didn't exit gracefully, force killing")
		if process.Process != nil {
			if err := process.Process.Kill(); err != nil {
				s.sendEvent("error", fmt.Sprintf("failed to kill process: %v", err))
			} else {
				s.sendEvent("update", "server process force killed")
			}
			process.Wait()
		}
	}
}

func (s *serverService) validatePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", fmt.Errorf("could not resolve path: %v", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("file does not exist: %v", err)
	}

	if !strings.HasSuffix(absPath, ".go") {
		return "", errors.New("file is not .go type")
	}

	return absPath, nil

}

func (s *serverService) sendEvent(eventType, message string) {
	// FIX: Thread-safe access to updateChan to prevent race conditions
	// Previous implementation accessed updateChan without mutex protection
	s.serverMu.Lock()
	updateChan := s.updateChan
	s.serverMu.Unlock()

	// Only send if channel is available
	if updateChan != nil {
		select {
		case updateChan <- UIEvent{Type: eventType, Message: message}:
		default:
		}
	}
}
