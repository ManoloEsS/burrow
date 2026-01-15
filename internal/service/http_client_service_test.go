package service

import (
	"testing"

	"github.com/ManoloEsS/burrow/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRequestJSONToStruct(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
		description string
	}{
		{
			name:        "Valid JSON request",
			jsonData:    "{\"name\":\"test\",\"method\":\"GET\",\"url\":\"http://example.com\",\"body\":\"test body\"}",
			expectError: false,
			description: "Should parse valid JSON request correctly",
		},
		{
			name:        "Invalid JSON",
			jsonData:    `{"name":"test","method":"GET","url":"http://example.com","body":"test body"`,
			expectError: true,
			description: "Should return error for invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := requestJSONToStruct([]byte(tt.jsonData))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, "test", req.Name)
				assert.Equal(t, "GET", req.Method)
				assert.Equal(t, "http://example.com", req.URL)
				assert.Equal(t, "test body", req.Body)
			}
		})
	}
}

func TestAddParams(t *testing.T) {
	tests := []struct {
		name           string
		params         map[string]string
		url            string
		expectedResult string
	}{
		{
			name: "URL with parameters",
			params: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			url:            "http://example.com",
			expectedResult: "http://example.com?param1=value1&param2=value2",
		},
		{
			name:           "Empty parameters",
			params:         map[string]string{},
			url:            "http://example.com",
			expectedResult: "http://example.com",
		},
		{
			name: "URL with single parameter",
			params: map[string]string{
				"param1": "value1",
			},
			url:            "http://example.com",
			expectedResult: "http://example.com?param1=value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addParams(tt.params, tt.url)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestReqStructToHttpReq(t *testing.T) {
	tests := []struct {
		name         string
		request      *domain.Request
		expectError  bool
		expectedSubs []string
		description  string
	}{
		{
			name: "Valid POST request with body",
			request: &domain.Request{
				Method: "POST",
				URL:    "http://example.com",
				Body:   "test body",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				ContentType: map[string]string{
					"Content-Type": "application/json",
				},
			},
			expectError:  false,
			expectedSubs: []string{"POST", "http://example.com", "application/json"},
		},
		{
			name: "GET request without body",
			request: &domain.Request{
				Method: "GET",
				URL:    "http://example.com",
				Body:   "",
			},
			expectError:  false,
			expectedSubs: []string{"GET", "http://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpReq, err := reqStructToHttpReq(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, httpReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, httpReq)

				for _, expectedSub := range tt.expectedSubs {
					if expectedSub == tt.request.Method {
						assert.Equal(t, tt.request.Method, httpReq.Method)
					} else if expectedSub == tt.request.URL {
						assert.Equal(t, tt.request.URL, httpReq.URL)
					} else {
						assert.Contains(t, httpReq.Header.Get("Content-Type"), expectedSub)
					}
				}
			}
		})
	}
}
