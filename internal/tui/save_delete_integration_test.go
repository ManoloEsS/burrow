package tui

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Test save and delete request list updates with proper synchronization
func TestSaveDeleteRequestListUpdates(t *testing.T) {
	cfg := &config.Config{DefaultPort: "8080"}
	tui := NewTui(cfg)
	tui.Components = createTuiLayout()

	// Mock service that tracks save/delete operations
	var savedRequests []*domain.Request
	mockService := &mockHttpServiceWithTracking{requests: &savedRequests}
	tui.HttpService = mockService

	// Test save functionality
	tui.Components.NameInput.SetText("test-request")
	tui.Components.MethodDropdown.SetCurrentOption(0) // GET
	tui.Components.URLInput.SetText("http://example.com")
	tui.Ui.SetFocus(tui.Components.Form) // Set proper focus

	// Get current request (should work)
	err := tui.getCurrentRequest()
	assert.NoError(t, err)
	assert.NotNil(t, tui.State.CurrentRequest)

	// Test save request
	tui.handleSaveRequest()

	// Verify request was saved
	assert.Equal(t, 1, len(savedRequests))
	assert.Equal(t, "test-request", savedRequests[0].Name)

	// Test delete functionality
	// Simulate having the request in the list
	tui.Components.RequestList.Clear()
	tui.Components.RequestList.AddItem("test-request      |  GET    |http://example.com", "", 0, nil)
	tui.State.CurrentFocused = tui.Components.RequestList

	// Test delete request
	tui.handleDeleteRequest()

	// Verify request was deleted
	assert.Equal(t, 0, len(savedRequests))
}

// Mock service that tracks save/delete operations for testing
type mockHttpServiceWithTracking struct {
	requests *[]*domain.Request
}

func (m *mockHttpServiceWithTracking) GetSavedRequests() ([]*domain.Request, error) {
	return *m.requests, nil
}

func (m *mockHttpServiceWithTracking) SendRequest(req *domain.Request) (*domain.Response, error) {
	return &domain.Response{Status: "200 OK"}, nil
}

func (m *mockHttpServiceWithTracking) SaveRequest(req *domain.Request) error {
	// Check if request already exists
	for _, existingReq := range *m.requests {
		if existingReq.Name == req.Name {
			return nil // Update existing (simplified for test)
		}
	}
	// Add new request
	*m.requests = append(*m.requests, req)
	return nil
}

func (m *mockHttpServiceWithTracking) DeleteRequest(name string) error {
	// Extract name from display format if needed
	if len(name) > 15 {
		name = name[:15] // Simplified extraction for test
	}

	// Remove request with matching name
	for i, req := range *m.requests {
		if req.Name == name {
			*m.requests = append((*m.requests)[:i], (*m.requests)[i+1:]...)
			break
		}
	}
	return nil
}
