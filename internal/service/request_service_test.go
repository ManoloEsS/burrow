package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequestService(t *testing.T) {
	tests := []struct {
		name     string
		db       *database.Database
		config   *config.Config
		callback RequestUpdateCallBack
	}{
		{
			name:     "valid request service creation",
			db:       &database.Database{},
			config:   &config.Config{DefaultPort: ":8080"},
			callback: func(resp *domain.Response) {},
		},
		{
			name:     "nil callback",
			db:       &database.Database{},
			config:   &config.Config{DefaultPort: ":8080"},
			callback: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewRequestService(tt.db, tt.config, tt.callback)
			assert.NotNil(t, service)

			rs, ok := service.(*requestService)
			require.True(t, ok)
			assert.Equal(t, tt.db, rs.requestRepo)
			assert.Equal(t, tt.config, rs.config)

			if tt.callback == nil {
				assert.Nil(t, rs.updateCallback)
			} else {
				assert.NotNil(t, rs.updateCallback)
			}
		})
	}
}

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "test response", "method": "` + r.Method + `"}`))
	}))
	defer server.Close()

	timeoutServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second) // Longer than client timeout (5s)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "this should timeout"}`))
	}))
	defer timeoutServer.Close()

	tests := []struct {
		name         string
		request      *domain.Request
		expectError  bool
		expectStatus int
		expectBody   string
	}{
		{
			name: "successful GET request",
			request: &domain.Request{
				Method: "GET",
				URL:    server.URL,
			},
			expectError:  false,
			expectStatus: 200,
			expectBody:   `{"message": "test response", "method": "GET"}`,
		},
		{
			name: "POST request with body",
			request: &domain.Request{
				Method:      "POST",
				URL:         server.URL,
				Body:        `{"test": "data"}`,
				ContentType: "application/json",
			},
			expectError:  false,
			expectStatus: 200,
			expectBody:   `{"message": "test response", "method": "POST"}`,
		},
		{
			name: "request with headers",
			request: &domain.Request{
				Method: "GET",
				URL:    server.URL,
				Headers: map[string]string{
					"Authorization": "Bearer token123",
					"Accept":        "application/json",
				},
			},
			expectError:  false,
			expectStatus: 200,
			expectBody:   `{"message": "test response", "method": "GET"}`,
		},
		{
			name: "invalid URL",
			request: &domain.Request{
				Method: "GET",
				URL:    "invalid-url",
			},
			expectError: true,
		},
		{
			name: "timeout scenario",
			request: &domain.Request{
				Method: "GET",
				URL:    timeoutServer.URL,
			},
			expectError: true,
		},
		{
			name: "empty request body",
			request: &domain.Request{
				Method: "POST",
				URL:    server.URL,
				Body:   "",
			},
			expectError:  false,
			expectStatus: 200,
			expectBody:   `{"message": "test response", "method": "POST"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &database.Database{}
			cfg := &config.Config{DefaultPort: ":8080"}
			callbackCalled := false
			var callbackResponse *domain.Response

			service := NewRequestService(mockDB, cfg, func(resp *domain.Response) {
				callbackCalled = true
				callbackResponse = resp
			})

			// Execute
			resp, err := service.SendRequest(tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.NotNil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, tt.expectStatus, resp.StatusCode)
			assert.Greater(t, resp.ResponseTime, time.Duration(0))
			assert.Equal(t, "application/json", resp.ContentType)
			assert.Equal(t, tt.expectBody, resp.Body)

			// Verify callback was called
			assert.True(t, callbackCalled)
			assert.NotNil(t, callbackResponse)
			assert.Equal(t, resp, callbackResponse)
		})
	}
}

func TestSaveRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.Request
		expectError bool
		description string
	}{
		{
			name: "save valid request",
			request: &domain.Request{
				Method:  "GET",
				URL:     "http://example.com",
				Headers: map[string]string{"Content-Type": "application/json"},
				Params:  map[string]string{"key": "value"},
				Body:    "test body",
			},
			expectError: false,
			description: "Should save request successfully when implemented",
		},
		{
			name: "save request with minimal data",
			request: &domain.Request{
				Method: "POST",
				URL:    "http://example.com",
			},
			expectError: false,
			description: "Should handle minimal requests",
		},
		{
			name:        "save nil request",
			request:     nil,
			expectError: false,
			description: "Should handle nil request gracefully",
		},
		{
			name: "save request with empty maps",
			request: &domain.Request{
				Method:  "PUT",
				URL:     "http://example.com",
				Headers: map[string]string{},
				Params:  map[string]string{},
				Body:    "",
			},
			expectError: false,
			description: "Should handle empty headers and params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &database.Database{}
			cfg := &config.Config{}
			service := NewRequestService(mockDB, cfg, nil)

			err := service.SaveRequest(tt.request)

			assert.NoError(t, err)

			// TODO: When implemented, add proper assertions
			t.Logf("TODO: Verify request was saved to database - %s", tt.description)
		})
	}
}

func TestGetSavedRequests(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
		description string
	}{
		{
			name:        "get saved requests",
			expectError: false,
			description: "Should return saved requests when implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &database.Database{}
			cfg := &config.Config{}
			service := NewRequestService(mockDB, cfg, nil)

			err := service.GetSavedRequests()

			assert.NoError(t, err)

			// TODO: When implemented, add proper assertions
			t.Logf("TODO: Verify requests are retrieved from database - %s", tt.description)
		})
	}
}

func TestResponseStringBuilder(t *testing.T) {
	tests := []struct {
		name     string
		request  *domain.Request
		contains []string
	}{
		{
			name: "GET request without body",
			request: &domain.Request{
				Method: "GET",
				URL:    "http://example.com",
			},
			contains: []string{
				"Mock Response for GET http://example.com",
				"Status: 200 OK",
				`{"message": "This is a mock response", "method": "GET"}`,
			},
		},
		{
			name: "POST request with body",
			request: &domain.Request{
				Method: "POST",
				URL:    "http://example.com",
				Body:   `{"test": "data"}`,
			},
			contains: []string{
				"Mock Response for POST http://example.com",
				"Echo Body:",
				`{"test": "data"}`,
			},
		},
		{
			name: "request with empty body",
			request: &domain.Request{
				Method: "PUT",
				URL:    "http://example.com",
				Body:   "",
			},
			contains: []string{
				"Mock Response for PUT http://example.com",
				`{"message": "This is a mock response", "method": "PUT"}`,
			},
		},
		{
			name: "DELETE request",
			request: &domain.Request{
				Method: "DELETE",
				URL:    "http://example.com/api/users/123",
			},
			contains: []string{
				"Mock Response for DELETE http://example.com/api/users/123",
				"Status: 200 OK",
				`{"message": "This is a mock response", "method": "DELETE"}`,
			},
		},
		{
			name: "request with whitespace body",
			request: &domain.Request{
				Method: "PATCH",
				URL:    "http://example.com",
				Body:   "   \n\t   ",
			},
			contains: []string{
				"Mock Response for PATCH http://example.com",
				"Echo Body:",
				"   \n\t   ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &requestService{}

			result := service.ResponseStringBuilder(tt.request)

			for _, expected := range tt.contains {
				assert.Contains(t, result, expected, "Response should contain: %s", expected)
			}
		})
	}
}

func TestMapToString(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		input := map[string]string{}
		result := mapToString(input)
		assert.Equal(t, "", result)
	})

	t.Run("nil map", func(t *testing.T) {
		var input map[string]string = nil
		result := mapToString(input)
		assert.Equal(t, "", result)
	})

	t.Run("single item", func(t *testing.T) {
		input := map[string]string{
			"key1": "value1",
		}
		result := mapToString(input)
		assert.Equal(t, "key1:value1", result)
	})

	t.Run("multiple items", func(t *testing.T) {
		input := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		result := mapToString(input)

		expectedParts := []string{"key1:value1", "key2:value2", "key3:value3"}
		for _, part := range expectedParts {
			assert.Contains(t, result, part)
		}

		resultParts := strings.Fields(result)
		assert.Equal(t, 3, len(resultParts))
	})

	t.Run("special characters", func(t *testing.T) {
		input := map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		}
		result := mapToString(input)

		assert.Contains(t, result, "Authorization:Bearer token123")
		assert.Contains(t, result, "Content-Type:application/json")

		resultParts := strings.Fields(result)
		assert.GreaterOrEqual(t, len(resultParts), 2)
	})

	t.Run("empty values", func(t *testing.T) {
		input := map[string]string{
			"empty":  "",
			"space":  " ",
			"normal": "value",
		}
		result := mapToString(input)

		assert.Contains(t, result, "empty:")
		assert.Contains(t, result, "space: ")
		assert.Contains(t, result, "normal:value")

		resultParts := strings.Fields(result)
		assert.Equal(t, 3, len(resultParts))
	})

	t.Run("keys and values with spaces", func(t *testing.T) {
		input := map[string]string{
			"User Agent":      "Mozilla/5.0",
			"Accept Encoding": "gzip, deflate",
		}
		result := mapToString(input)

		assert.Contains(t, result, "User Agent:Mozilla/5.0")
		assert.Contains(t, result, "Accept Encoding:gzip, deflate")

		resultParts := strings.Fields(result)
		assert.GreaterOrEqual(t, len(resultParts), 2)
	})
}

func TestRequestServiceIntegration(t *testing.T) {
	t.Run("callback functionality", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}))
		defer server.Close()

		callbackResponses := make([]*domain.Response, 0)
		cfg := &config.Config{DefaultPort: ":8080"}
		service := NewRequestService(&database.Database{}, cfg, func(resp *domain.Response) {
			callbackResponses = append(callbackResponses, resp)
		})

		request := &domain.Request{
			Method: "GET",
			URL:    server.URL,
		}

		for i := 0; i < 3; i++ {
			resp, err := service.SendRequest(request)
			require.NoError(t, err)
			require.NotNil(t, resp)
		}

		assert.Len(t, callbackResponses, 3)

		for _, resp := range callbackResponses {
			assert.Equal(t, 200, resp.StatusCode)
			assert.Contains(t, resp.Body, `"status": "ok"`)
		}
	})
}

func BenchmarkSendRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "benchmark test"}`))
	}))
	defer server.Close()

	cfg := &config.Config{DefaultPort: ":8080"}
	service := NewRequestService(&database.Database{}, cfg, nil)
	request := &domain.Request{
		Method: "GET",
		URL:    server.URL,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.SendRequest(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}
