package tui

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/stretchr/testify/assert"
)

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
			tui.Components = createTuiLayout()

			tui.Components.NameInput.SetText("test-request")
			tui.Components.MethodDropdown.SetCurrentOption(0)
			tui.Components.URLInput.SetText("http://example.com")

			if tt.name == "Get request with all fields" {
				tui.Components.HeadersText.SetText("Content-Type:application/json", true)
				tui.Components.ParamsText.SetText("param1:value1", true)
				tui.Components.BodyText.SetText("test body", true)
				tui.Components.BodyType.SetCurrentOption(1)
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

func TestMapToString(t *testing.T) {
	tests := []struct {
		name         string
		inputMap     map[string]string
		expectedSubs []string
		description  string
	}{
		{
			name: "Multiple headers",
			inputMap: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token",
			},
			expectedSubs: []string{"Content-Type:application/json", "Authorization:Bearer token"},
			description:  "Should convert map with multiple items to string",
		},
		{
			name:         "Empty map",
			inputMap:     map[string]string{},
			expectedSubs: []string{""},
			description:  "Should return empty string for empty map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToString(tt.inputMap)

			for _, expectedSub := range tt.expectedSubs {
				if expectedSub == "" {
					assert.Equal(t, "", result)
				} else {
					assert.Contains(t, result, expectedSub)
				}
			}
		})
	}
}
