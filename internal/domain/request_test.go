package domain

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Parse GET method",
			input:       "get",
			expected:    "GET",
			description: "Should convert lowercase GET to uppercase",
		},
		{
			name:        "Parse POST method",
			input:       "post",
			expected:    "POST",
			description: "Should convert lowercase POST to uppercase",
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
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Parse domain without protocol",
			input:       "example.com",
			expected:    "https://example.com",
			description: "Should add HTTPS protocol to domain",
		},
		{
			name:        "Parse localhost with port",
			input:       "localhost:3000",
			expected:    "https://localhost:3000",
			description: "Should preserve port specification",
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
	}{
		{
			name:        "Parse JSON headers",
			input:       "Content-Type:application/json Authorization:Bearer:token",
			description: "Should parse multiple headers with colons",
		},
		{
			name:        "Parse simple headers",
			input:       "Accept:text/plain",
			description: "Should parse headers without special characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseHeaders(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, "Burrow/1.0 (+https://github.com/ManoloEsS/burrow)", req.Headers["User-Agent"])
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
		description string
	}{
		{
			name:        "Parse JSON body",
			input:       `{"key": "value"}`,
			description: "Should parse JSON body content",
		},
		{
			name:        "Parse text body",
			input:       "simple text body",
			description: "Should parse plain text body content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseBody(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, req.Body)
		})
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Parse URL parameters",
			input:       "param1:value1 param2:value2",
			description: "Should parse multiple parameters",
		},
		{
			name:        "Parse API key parameters",
			input:       "api_key:12345 limit:25",
			description: "Should parse API configuration parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			err := req.ParseParams(tt.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, req.Params)
		})
	}
}

func TestParseName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Parse camel case name",
			input:       "MyRequest",
			expected:    "myrequest",
			description: "Should convert camel case to lowercase",
		},
		{
			name:        "Parse spaced name",
			input:       "Test Request Name",
			expected:    "test request name",
			description: "Should preserve spaces and lowercase",
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "Build JSON request",
			description: "Should build request with JSON body type",
		},
		{
			name:        "Build form request",
			description: "Should build request with form data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest()
			cfg := &config.Config{DefaultPort: "8080"}
			bodyType := "JSON"
			if tt.name == "Build form request" {
				bodyType = "FORM"
			}
			err := req.BuildRequest("test", "GET", "example.com", "Accept:application/json", "param:value", bodyType, "test body", cfg)
			assert.NoError(t, err)
			assert.Equal(t, "test", req.Name)
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "https://example.com", req.URL)
			assert.NotNil(t, req.ContentType["Content-Type"])
		})
	}
}
