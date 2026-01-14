package domain

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildResponse(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		statusCode     int
		contentType    string
		body           string
		expectedStatus string
		expectedType   string
		description    string
	}{
		{
			name:           "Build JSON response",
			status:         "200 OK",
			statusCode:     200,
			contentType:    "application/json",
			body:           `{"message": "success"}`,
			expectedStatus: "200 OK",
			expectedType:   "application/json",
			description:    "Should build response with JSON content",
		},
		{
			name:           "Build HTML response",
			status:         "404 Not Found",
			statusCode:     404,
			contentType:    "text/html",
			body:           "<h1>Not Found</h1>",
			expectedStatus: "404 Not Found",
			expectedType:   "text/html",
			description:    "Should build response with HTML content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP response
			body := io.NopCloser(strings.NewReader(tt.body))
			mockHttpResponse := &http.Response{
				Status:        tt.status,
				StatusCode:    tt.statusCode,
				Header:        make(http.Header),
				Body:          body,
				ContentLength: int64(len(tt.body)),
			}
			mockHttpResponse.Header.Set("Content-Type", tt.contentType)

			// Test the BuildResponse function
			resp := &Response{}
			err := resp.BuildResponse(mockHttpResponse)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.Status)
			assert.Equal(t, tt.expectedType, resp.ContentType)
			assert.Equal(t, int64(len(tt.body)), resp.ContentLenght)
			assert.Equal(t, tt.body, resp.Body)
		})
	}
}
