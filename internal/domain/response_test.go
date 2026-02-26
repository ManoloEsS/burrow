package domain

import (
	"bytes"
	"encoding/json"
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
	}{
		{
			name:           "Build JSON response",
			status:         "200 OK",
			statusCode:     200,
			contentType:    "application/json",
			body:           "{\n \"message\": \"success\"\n}",
			expectedStatus: "200 OK",
			expectedType:   "application/json",
		},
		{
			name:           "Build text response",
			status:         "404 Not Found",
			statusCode:     404,
			contentType:    "text/html",
			body:           "<h1>Not Found</h1>",
			expectedStatus: "404 Not Found",
			expectedType:   "text/html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.body))
			mockHttpResponse := &http.Response{
				Status:        tt.status,
				StatusCode:    tt.statusCode,
				Header:        make(http.Header),
				Body:          body,
				ContentLength: int64(len(tt.body)),
			}
			mockHttpResponse.Header.Set("Content-Type", tt.contentType)

			resp := &Response{}
			err := resp.BuildResponse(mockHttpResponse)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.Status)
			assert.Equal(t, tt.expectedType, resp.ContentType)
			assert.Equal(t, int64(len(tt.body)), resp.ContentLenght)
			if strings.HasPrefix(tt.contentType, "application/json") {
				var prettyJson bytes.Buffer
				if err := json.Indent(&prettyJson, []byte(tt.body), "", " "); err != nil {
					assert.Equal(t, tt.body, resp.Body)
				} else {
					assert.Equal(t, prettyJson.String(), resp.Body)
				}
			} else {
				assert.Equal(t, tt.body, resp.Body)
			}
		})
	}
}
