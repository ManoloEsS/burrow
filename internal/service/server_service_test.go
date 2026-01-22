package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartServerInvalidPath(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	updateChan := make(chan UIEvent, 10)
	defer close(updateChan)

	err := service.StartServer("/nonexistent/file.go", "8080", updateChan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
	assert.Empty(t, serverService.pathToServer)
	assert.False(t, serverService.isRunning)
}

func TestStartServerNonGoFile(t *testing.T) {
	service := NewServerService()

	tempDir := t.TempDir()
	serverFile := filepath.Join(tempDir, "test_server.txt")
	err := os.WriteFile(serverFile, []byte("not a go file"), 0644)
	require.NoError(t, err)

	updateChan := make(chan UIEvent, 10)
	defer close(updateChan)

	err = service.StartServer(serverFile, "8080", updateChan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file is not .go type")
}

func TestStopServer(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	serverService.isRunning = true
	serverService.cancelFunc = func() {}
	serverService.serverProcess = nil

	updateChan := make(chan UIEvent, 10)
	defer close(updateChan)
	serverService.updateChan = updateChan

	err := service.StopServer()

	assert.NoError(t, err)
	assert.Nil(t, serverService.cancelFunc)
}

func TestStopServerNotRunning(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	assert.False(t, serverService.isRunning)
	assert.Nil(t, serverService.cancelFunc)

	err := service.StopServer()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server not running")
}

func TestValidatePath(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	tempDir := t.TempDir()
	serverFile := filepath.Join(tempDir, "test_server.go")
	err := os.WriteFile(serverFile, []byte("package main\n\nfunc main() {}"), 0644)
	require.NoError(t, err)

	validPath, err := serverService.validatePath(serverFile)

	assert.NoError(t, err)
	assert.Equal(t, serverFile, validPath)

	invalidPath, err := serverService.validatePath("/nonexistent/file.go")

	assert.Error(t, err)
	assert.Empty(t, invalidPath)
	assert.Contains(t, err.Error(), "file does not exist")

	nonGoFile := filepath.Join(tempDir, "test_server.txt")
	err = os.WriteFile(nonGoFile, []byte("not go"), 0644)
	require.NoError(t, err)

	invalidPath2, err := serverService.validatePath(nonGoFile)

	assert.Error(t, err)
	assert.Empty(t, invalidPath2)
	assert.Contains(t, err.Error(), "file is not .go type")
}

func TestCleanupBinaryNoBinary(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	assert.Empty(t, serverService.binaryPath)

	serverService.cleanupBinary()

	assert.Empty(t, serverService.binaryPath)
}

func TestSendEvent(t *testing.T) {
	service := NewServerService()
	serverService := service.(*serverService)

	eventChan := make(chan UIEvent, 10)
	defer close(eventChan)
	serverService.updateChan = eventChan

	serverService.sendEvent("test", "test message")

	select {
	case event := <-eventChan:
		assert.Equal(t, "test", event.Type)
		assert.Equal(t, "test message", event.Message)
	default:
	}
}
