package domain

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildResponse tests the BuildResponse function
func TestBuildResponse(t *testing.T) {
	tests := []struct {
		name           string
		httpResponse   *http.Response
		expectedStatus int
		expectedBody   string
		expectedType   string
		expectedError  bool
		description    string
	}{
		{
			name:           "successful JSON response",
			httpResponse:   createTestResponse(http.StatusOK, "application/json", `{"message": "success"}`),
			expectedStatus: 200,
			expectedBody:   `{"message": "success"}`,
			expectedType:   "application/json",
			expectedError:  false,
			description:    "Should build response successfully with JSON",
		},
		{
			name:           "HTML response",
			httpResponse:   createTestResponse(http.StatusOK, "text/html", "<html><body>Test</body></html>"),
			expectedStatus: 200,
			expectedBody:   "<html><body>Test</body></html>",
			expectedType:   "text/html",
			expectedError:  false,
			description:    "Should handle HTML responses",
		},
		{
			name:           "empty response body",
			httpResponse:   createTestResponse(http.StatusNoContent, "application/json", ""),
			expectedStatus: 204,
			expectedBody:   "",
			expectedType:   "application/json",
			expectedError:  false,
			description:    "Should handle empty response body",
		},
		{
			name:           "large response body",
			httpResponse:   createTestResponse(http.StatusOK, "text/plain", strings.Repeat("test", 1000)),
			expectedStatus: 200,
			expectedBody:   strings.Repeat("test", 1000),
			expectedType:   "text/plain",
			expectedError:  false,
			description:    "Should handle large response bodies",
		},
		{
			name: "response with custom headers",
			httpResponse: createCustomTestResponse(
				http.StatusOK,
				"application/xml",
				"<root>test</root>",
				map[string]string{
					"X-Custom-Header": "custom-value",
					"Cache-Control":   "no-cache",
				},
			),
			expectedStatus: 200,
			expectedBody:   "<root>test</root>",
			expectedType:   "application/xml",
			expectedError:  false,
			description:    "Should extract content-type from headers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{}

			err := resp.BuildResponse(tt.httpResponse)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedType, resp.ContentType)
			assert.Equal(t, tt.expectedBody, resp.Body)

			// ContentLength should be set correctly
			if tt.httpResponse.ContentLength > 0 {
				assert.Equal(t, int(tt.httpResponse.ContentLength), resp.ContentLenght)
			}
		})
	}
}

// TestResponseStruct tests the Response struct fields
func TestResponseStruct(t *testing.T) {
	resp := &Response{}

	// Test initial state
	assert.Equal(t, 0, resp.StatusCode)
	assert.Equal(t, "", resp.ContentType)
	assert.Equal(t, 0, resp.ContentLenght)
	assert.Equal(t, "", resp.Body)
	assert.Equal(t, time.Duration(0), resp.ResponseTime)

	// Test field assignments
	resp.StatusCode = 200
	resp.ContentType = "application/json"
	resp.ContentLenght = 100
	resp.Body = `{"test": "data"}`
	resp.ResponseTime = time.Millisecond * 150

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.ContentType)
	assert.Equal(t, 100, resp.ContentLenght)
	assert.Equal(t, `{"test": "data"}`, resp.Body)
	assert.Equal(t, time.Millisecond*150, resp.ResponseTime)
}

// TestBuildResponseEdgeCases tests edge cases and error conditions
func TestBuildResponseEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *http.Response
		description string
		expectError bool
	}{
		{
			name: "response with negative content length",
			setupFunc: func() *http.Response {
				resp := createTestResponse(http.StatusOK, "text/plain", "test")
				resp.ContentLength = -1 // Unknown content length
				return resp
			},
			description: "Should handle unknown content length",
			expectError: false,
		},
		{
			name: "response with no content-type header",
			setupFunc: func() *http.Response {
				resp := httptest.NewRecorder()
				resp.WriteHeader(http.StatusOK)
				resp.WriteString("test")
				return resp.Result()
			},
			description: "Should handle missing content-type",
			expectError: false,
		},
		{
			name: "response with multiple content-type values",
			setupFunc: func() *http.Response {
				rec := httptest.NewRecorder()
				rec.Header().Set("Content-Type", "application/json; charset=utf-8")
				rec.WriteHeader(http.StatusOK)
				rec.WriteString(`{"test": "data"}`)
				return rec.Result()
			},
			description: "Should handle complex content-type headers",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpResp := tt.setupFunc()
			resp := &Response{}

			err := resp.BuildResponse(httpResp)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, httpResp.StatusCode, resp.StatusCode)

			// The body should contain something (even if error message)
			assert.NotEmpty(t, resp.Body)
		})
	}
}

// TestBuildResponseErrorHandling tests error handling in BuildResponse
func TestBuildResponseErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *http.Response
		description string
		expectError bool
	}{
		{
			name: "nil body reader",
			setupFunc: func() *http.Response {
				rec := httptest.NewRecorder()
				rec.WriteHeader(http.StatusOK)
				result := rec.Result()
				result.Body = nil // Explicitly set to nil
				return result
			},
			description: "Should handle nil body gracefully",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpResp := tt.setupFunc()
			resp := &Response{}

			err := resp.BuildResponse(httpResp)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, resp.Body, "Error reading body:")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBuildResponseStatusCodes tests various HTTP status codes
func TestBuildResponseStatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}

	for _, status := range statusCodes {
		t.Run(http.StatusText(status), func(t *testing.T) {
			httpResp := createTestResponse(status, "application/json", `{"status": "`+http.StatusText(status)+`"}`)
			resp := &Response{}

			err := resp.BuildResponse(httpResp)

			assert.NoError(t, err)
			assert.Equal(t, status, resp.StatusCode)
			assert.Equal(t, "application/json", resp.ContentType)
			assert.Contains(t, resp.Body, http.StatusText(status))
		})
	}
}

// TestBuildResponseBodyContent tests various body content types
func TestBuildResponseBodyContent(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "JSON object",
			body:     `{"key": "value", "number": 123}`,
			expected: `{"key": "value", "number": 123}`,
		},
		{
			name:     "XML content",
			body:     `<root><item>test</item></root>`,
			expected: `<root><item>test</item></root>`,
		},
		{
			name:     "plain text",
			body:     "This is plain text content",
			expected: "This is plain text content",
		},
		{
			name:     "binary data as string",
			body:     string([]byte{0x00, 0x01, 0x02, 0x03}),
			expected: string([]byte{0x00, 0x01, 0x02, 0x03}),
		},
		{
			name:     "unicode content",
			body:     "Hello 世界 🌍",
			expected: "Hello 世界 🌍",
		},
		{
			name:     "whitespace content",
			body:     "   \n\t   \n  ",
			expected: "   \n\t   \n  ",
		},
		{
			name:     "empty string",
			body:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpResp := createTestResponse(http.StatusOK, "text/plain", tt.body)
			resp := &Response{}

			err := resp.BuildResponse(httpResp)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, resp.Body)
		})
	}
}

// Performance test for response building
func TestBuildResponsePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	largeBody := strings.Repeat("data", 10000) // 40KB
	httpResp := createTestResponse(http.StatusOK, "application/json", largeBody)

	start := time.Now()
	for i := 0; i < 1000; i++ {
		resp := &Response{}
		resp.BuildResponse(httpResp)
	}
	duration := time.Since(start)

	t.Logf("Built 1000 responses in %v (%.2fms per response)",
		duration, float64(duration.Nanoseconds())/1e6/1000)

	// Should complete reasonably quickly (adjust threshold as needed)
	assert.Less(t, duration, time.Second)
}

// Benchmark test
func BenchmarkBuildResponse(b *testing.B) {
	httpResp := createTestResponse(http.StatusOK, "application/json", `{"test": "data"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &Response{}
		resp.BuildResponse(httpResp)
	}
}

// Helper functions for testing
func createTestResponse(status int, contentType, body string) *http.Response {
	rec := httptest.NewRecorder()
	if contentType != "" {
		rec.Header().Set("Content-Type", contentType)
	}
	rec.WriteHeader(status)
	rec.WriteString(body)
	return rec.Result()
}

func createCustomTestResponse(status int, contentType, body string, headers map[string]string) *http.Response {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", contentType)

	for key, value := range headers {
		rec.Header().Set(key, value)
	}

	rec.WriteHeader(status)
	rec.WriteString(body)
	return rec.Result()
}
