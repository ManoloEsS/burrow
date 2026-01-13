package domain

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/stretchr/testify/assert"
)

//TODO: tests for all parsing methods

func TestParseMethod(t *testing.T) {
	reqParseMethodTests := []struct {
		name         string
		input        string
		expectError  bool
		expectResult string
	}{
		{
			name:         "test lowercase method",
			input:        "get",
			expectError:  false,
			expectResult: "GET",
		},
		{
			name:         "test uppercase method",
			input:        "GET",
			expectError:  false,
			expectResult: "GET",
		},
		{
			name:         "test uppercase and lowercase method",
			input:        "GeT",
			expectError:  false,
			expectResult: "GET",
		},
		{
			name:         "test empty string",
			input:        "",
			expectError:  true,
			expectResult: "",
		},
	}

	for _, tt := range reqParseMethodTests {
		t.Run(tt.name, func(t *testing.T) {
			request := Request{}
			err := request.ParseMethod(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, request.Method, tt.expectResult)
		})
	}
}

func TestParseUrl(t *testing.T) {
	reqParseUrlTests := []struct {
		name         string
		input        string
		expectError  bool
		expectResult string
	}{
		{
			name:         "test no url",
			input:        "",
			expectError:  false,
			expectResult: "http://localhost:8080",
		},
		{
			name:         "test default url",
			input:        "http://localhost:8080",
			expectError:  false,
			expectResult: "http://localhost:8080",
		},
		{
			name:         "test https:// prefix",
			input:        "https://someurl.com",
			expectError:  false,
			expectResult: "https://someurl.com",
		},
		{
			name:         "test http:// prefix",
			input:        "http://someurl.com",
			expectError:  false,
			expectResult: "http://someurl.com",
		},
		{
			name:         "test no prefix",
			input:        "someurl.com",
			expectError:  false,
			expectResult: "https://someurl.com",
		},
		{
			name:         "test default port",
			input:        ":8080",
			expectError:  false,
			expectResult: "http://localhost:8080",
		},
	}

	for _, tt := range reqParseUrlTests {
		t.Run(tt.name, func(t *testing.T) {
			request := Request{}
			cfg := config.Config{
				DefaultPort: ":8080",
			}
			err := request.ParseUrl(&cfg, tt.input)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, request.URL, tt.expectResult)
		})
	}
}
