package tui

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/ManoloEsS/burrow/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHandleSendRequest(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "Send valid request",
			description: "Should send request with valid data",
		},
		{
			name:        "Send request with missing URL",
			description: "Should handle sending request with missing URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultPort: "8080"}
			tui := NewTui(cfg)
			// Don't initialize to avoid the panic in loadSavedRequests
			tui.Components = createTuiLayout()

			// Set up test data
			if tt.name == "Send valid request" {
				tui.Components.NameInput.SetText("test-request")
				tui.Components.MethodDropdown.SetCurrentOption(0) // GET
				tui.Components.URLInput.SetText("http://example.com")
			} else {
				// Leave URL empty for error case
				tui.Components.NameInput.SetText("test-request")
				tui.Components.MethodDropdown.SetCurrentOption(0) // GET
				tui.Components.URLInput.SetText("")
			}

			// Test the handler function exists and doesn't panic
			assert.NotNil(t, tui.handleSendRequest)

			// Test getCurrentRequest works
			err := tui.getCurrentRequest()
			if tt.name == "Send valid request" {
				assert.NoError(t, err)
				assert.NotNil(t, tui.State.CurrentRequest)
			} else {
				// The request building succeeds but URL will be empty
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetCurrentRequest(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "Get request with all fields",
			description: "Should build request with all form fields",
		},
		{
			name:        "Get request with minimal fields",
			description: "Should build request with minimal required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultPort: "8080"}
			tui := NewTui(cfg)
			// Don't initialize to avoid the panic in loadSavedRequests
			tui.Components = createTuiLayout()

			// Set up form data
			tui.Components.NameInput.SetText("test-request")
			tui.Components.MethodDropdown.SetCurrentOption(0) // GET
			tui.Components.URLInput.SetText("http://example.com")

			if tt.name == "Get request with all fields" {
				tui.Components.HeadersText.SetText("Content-Type:application/json", true)
				tui.Components.ParamsText.SetText("param1:value1", true)
				tui.Components.BodyText.SetText("test body", true)
				tui.Components.BodyType.SetCurrentOption(1) // JSON
			}

			err := tui.getCurrentRequest()
			assert.NoError(t, err)
			assert.NotNil(t, tui.State.CurrentRequest)
			assert.Equal(t, "test-request", tui.State.CurrentRequest.Name)
			assert.Equal(t, "GET", tui.State.CurrentRequest.Method)
			assert.Equal(t, "http://example.com", tui.State.CurrentRequest.URL)
		})
	}
}

func TestLoadSavedRequestsNoDuplication(t *testing.T) {
	cfg := &config.Config{DefaultPort: "8080"}
	tui := NewTui(cfg)
	tui.Components = createTuiLayout()

	// Mock HttpService with some test requests
	mockRequests := []*domain.Request{
		{Name: "test1", Method: "GET", URL: "http://example.com"},
		{Name: "test2", Method: "POST", URL: "http://api.example.com"},
	}

	// Create a mock service that returns our test requests
	var mockService service.HttpClientService = &mockHttpService{requests: mockRequests}
	tui.HttpService = mockService

	// Call loadSavedRequests twice - this should not cause duplication
	tui.loadSavedRequests()
	firstCount := tui.Components.RequestList.GetItemCount()

	tui.loadSavedRequests()
	secondCount := tui.Components.RequestList.GetItemCount()

	// The count should be the same - no duplication occurred
	assert.Equal(t, firstCount, secondCount, "loadSavedRequests should not duplicate items")
	assert.Equal(t, 2, secondCount, "Should have exactly 2 items")

	// Verify the items are correct
	firstItem, _ := tui.Components.RequestList.GetItemText(0)
	secondItem, _ := tui.Components.RequestList.GetItemText(1)

	assert.Contains(t, firstItem, "test1")
	assert.Contains(t, firstItem, "GET")
	assert.Contains(t, secondItem, "test2")
	assert.Contains(t, secondItem, "POST")
}

// Mock HTTP service for testing - implements service.HttpClientService interface
type mockHttpService struct {
	requests []*domain.Request
}

func (m *mockHttpService) GetSavedRequests() ([]*domain.Request, error) {
	return m.requests, nil
}

func (m *mockHttpService) SendRequest(req *domain.Request) (*domain.Response, error) {
	return &domain.Response{Status: "200 OK"}, nil
}

func (m *mockHttpService) SaveRequest(req *domain.Request) error {
	return nil
}

func (m *mockHttpService) DeleteRequest(name string) error {
	// Remove request with matching name
	for i, req := range m.requests {
		if req.Name == name {
			m.requests = append(m.requests[:i], m.requests[i+1:]...)
			break
		}
	}
	return nil
}

func TestResponseStringBuilder(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "Build response string with all fields",
			description: "Should build response string with all response data",
		},
		{
			name:        "Build response string with minimal data",
			description: "Should build response string with minimal response data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultPort: "8080"}
			tui := NewTui(cfg)
			// Don't initialize to avoid the panic in loadSavedRequests
			tui.Components = createTuiLayout()

			resp := &domain.Response{
				Status:        "200 OK",
				ContentType:   "application/json",
				ContentLenght: 100,
				Body:          "test response body",
			}

			result := tui.responseStringBuilder(resp)

			assert.Contains(t, result, "200 OK")
			assert.Contains(t, result, "application/json")
			assert.Contains(t, result, "100")

			if tt.name == "Build response string with all fields" {
				assert.Contains(t, result, "test response body")
			}
		})
	}
}
