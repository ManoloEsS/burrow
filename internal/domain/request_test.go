package domain

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Parse GET method",
			input:    "get",
			expected: "GET",
		},
		{
			name:     "Parse POST method",
			input:    "post",
			expected: "POST",
		},
		{
			name:     "Parse PUT method",
			input:    "put",
			expected: "PUT",
		},
		{
			name:     "Parse DELETE method",
			input:    "delete",
			expected: "DELETE",
		},
		{
			name:     "Parse HEAD method",
			input:    "head",
			expected: "HEAD",
		},
		{
			name:     "Parse undefined method",
			input:    "undefined",
			expected: "UNDEFINED",
		},
		{
			name:     "Parse GET uppercase method",
			input:    "GET",
			expected: "GET",
		},
		{
			name:     "Parse GET method with trailing space",
			input:    "  get  ",
			expected: "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseMethod(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.Method)
		})
	}
}

func TestParseUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Parse domain without protocol",
			input:    "example.com",
			expected: "https://example.com",
		},
		{
			name:     "Parse localhost with port",
			input:    "localhost:3000",
			expected: "http://localhost:3000",
		},
		{
			name:     "Parse domain with http protocol",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "Parse domain with https protocol",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "Parse port",
			input:    ":3000",
			expected: "http://localhost:3000",
		},
		{
			name:     "Parse empty string",
			input:    "",
			expected: "http://localhost:8080",
		},
		{
			name:     "Parse invalid port syntax",
			input:    "3000",
			expected: "https://3000",
		},
		{
			name:     "Parse prefix localhost:",
			input:    "localhost:3000",
			expected: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			cfg := &config.Config{DefaultPort: "8080"}
			err := req.ParseUrl(cfg, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.URL)
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
		expected    map[string]string
	}{
		{
			name:     "Parse JSON headers",
			input:    "Content-Type:application/json, Authorization:Bearer token",
			expected: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token", "User-Agent": "Burrow/1.0.0(github.com/ManoloEsS/burrow)"},
		},
		{
			name:     "Parse simple headers",
			input:    "Accept:text/plain",
			expected: map[string]string{"Accept": "text/plain", "User-Agent": "Burrow/1.0.0(github.com/ManoloEsS/burrow)"},
		},
		{
			name:     "Parse simple headers trailing comma",
			input:    "Accept:text/plain,",
			expected: map[string]string{"Accept": "text/plain", "User-Agent": "Burrow/1.0.0(github.com/ManoloEsS/burrow)"},
		},
		{
			name:     "Parse simple headers trailing comma and space",
			input:    "Accept:text/plain, ",
			expected: map[string]string{"Accept": "text/plain", "User-Agent": "Burrow/1.0.0(github.com/ManoloEsS/burrow)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseHeaders(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.Headers)

		})
	}
}

func TestParseBodyType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Parse JSON body type",
			input:       "JSON",
			expected:    "application/json",
			description: "Should convert JSON to application/json",
		},
		{
			name:        "Parse Text body type",
			input:       "Text",
			expected:    "text/plain; charset=utf-8",
			description: "Should convert Text to plain text with charset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseBodyType(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.ContentType["Content-Type"])
		})
	}
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		bodyType    string
		expectError bool
	}{
		{
			name:        "Parse JSON body",
			input:       "{\"key\": \"value\"}",
			bodyType:    "JSON",
			expectError: false,
		},
		{
			name:        "Parse text body",
			input:       "simple text body",
			bodyType:    "Text",
			expectError: false,
		},
		{
			name:        "Parse invalid JSON body",
			input:       "{\"key\": \"value}",
			bodyType:    "JSON",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseBody(tt.input, tt.bodyType)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.input, req.Body)

			}
		})
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Parse single param",
			input:    "param1:value1",
			expected: map[string]string{"param1": "value1"},
		},
		{
			name:     "Parse multiple params",
			input:    "param1:value1, param2:value2",
			expected: map[string]string{"param1": "value1", "param2": "value2"},
		},
		{
			name:     "Parse simple params trailing comma",
			input:    "param1:value1, param2:value2,",
			expected: map[string]string{"param1": "value1", "param2": "value2"},
		},
		{
			name:     "Parse simple params trailing comma and space",
			input:    "param1:value1, param2:value2,  ",
			expected: map[string]string{"param1": "value1", "param2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseParams(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.Params)

		})
	}
}

func TestParseName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Parse camel case name",
			input:    "MyRequest",
			expected: "myrequest",
		},
		{
			name:     "Parse spaced name",
			input:    "Test Request Name",
			expected: "test request name",
		},
		{
			name:     "Parse empty name",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseName(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, req.Name)
		})
	}
}
