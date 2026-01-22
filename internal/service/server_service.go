package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ManoloEsS/burrow/internal/config"
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
	binaryPath     string
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
	s.updateChan = updateChan
	s.serverMu.Unlock()

	s.sendEvent("update", "starting server...")

	validPath, err := s.validatePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}
	s.sendEvent("update", "valid path")

	if s.isRunning {
		return fmt.Errorf("server already running")
	}

	s.serverMu.Lock()
	s.pathToServer = validPath
	s.healthCheckURL = "http://localhost:" + port + "/health"
	s.serverMu.Unlock()

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

	timeoutTicker := time.NewTicker(time.Minute * 15)

	for {
		select {
		case <-ctx.Done():
			healthCheckerCancel()
			wg.Wait()
			s.gracefulShutdown()

			s.serverMu.Lock()
			defer s.serverMu.Unlock()

			s.isRunning = false
			s.serverProcess = nil
			return

		case <-timeoutTicker.C:
			healthCheckerCancel()
			wg.Wait()
			s.gracefulShutdown()

			s.serverMu.Lock()
			defer s.serverMu.Unlock()

			s.isRunning = false
			s.serverProcess = nil
			return
		}
	}

}

func (s *serverService) buildBinary(path string) error {
	s.sendEvent("update", "building binary...")

	cacheDir := config.GetServerCachePath()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	binaryName := fmt.Sprintf("burrow-server-%x", md5.Sum([]byte(path)))[:8]
	binaryPath := filepath.Join(cacheDir, binaryName)

	cmd := exec.Command("go", "build", "-o", binaryPath, "-trimpath", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	s.binaryPath = binaryPath
	s.sendEvent("update", "server running...")
	return nil
}

func (s *serverService) runCmdFromPath(ctx context.Context) (*exec.Cmd, error) {
	if err := s.buildBinary(s.pathToServer); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, s.binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *serverService) healthChecker(ctx context.Context) {
	time.Sleep(time.Second * 1)

	s.sendEvent("update", "trying to reach server")

	resp, err := s.httpClient.Get(s.healthCheckURL)
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
		s.sendEvent("update", "server reached, starting health checker")
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := s.httpClient.Get(s.healthCheckURL)
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
	if s.cancelFunc == nil || !s.isRunning {
		return errors.New("server not running")
	}

	cancel := s.cancelFunc

	s.serverMu.Lock()
	defer s.serverMu.Unlock()

	s.cancelFunc = nil
	cancel()

	return nil
}

func (s *serverService) gracefulShutdown() {
	if s.serverProcess == nil || s.serverProcess.Process == nil {
		s.sendEvent("update", "no server process to stop")
		s.cleanupBinary()
		return
	}

	s.sendEvent("update", "stopping server")

	process := s.serverProcess

	if err := process.Process.Signal(syscall.SIGTERM); err != nil {
		s.sendEvent("error", fmt.Sprintf("failed to terminate process: %v", err))
	}

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
		s.sendEvent("error", "server didn't shutdown gracefully, force killing")
		if process.Process != nil {
			if err := process.Process.Kill(); err != nil {
				s.sendEvent("error", fmt.Sprintf("failed to kill process %d: %v", s.serverProcess.Process.Pid, err))
			} else {
				s.sendEvent("update", "server process force killed")
			}
			process.Wait()
		}
	}

	s.cleanupBinary()

	s.sendEvent("update", "server not running...ready")
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

func (s *serverService) cleanupBinary() {
	if s.binaryPath != "" {
		if err := os.Remove(s.binaryPath); err != nil && !os.IsNotExist(err) {
			s.sendEvent("error", fmt.Sprintf("failed to cleanup binary: %v", err))
		} else {
			s.sendEvent("update", "cleanup successful")
		}
		s.binaryPath = ""
	}
}

func (s *serverService) sendEvent(eventType, message string) {
	s.serverMu.Lock()
	updateChan := s.updateChan
	s.serverMu.Unlock()

	if updateChan != nil {
		select {
		case updateChan <- UIEvent{Type: eventType, Message: message}:
		default:
		}
	}
}
