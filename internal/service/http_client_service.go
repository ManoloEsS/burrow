package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
)

type httpClientService struct {
	requestRepo *database.Database
	config      *config.Config
}

func NewHttpClientService(requestRepo *database.Database, config *config.Config) HttpClientService {
	return &httpClientService{
		requestRepo: requestRepo,
		config:      config,
	}
}

func (s *httpClientService) SendRequest(req *domain.Request) (*domain.Response, error) {
	// Validate request to prevent panics and provide better error messages
	if req.Method == "" {
		return nil, fmt.Errorf("HTTP method is required")
	}
	if req.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	if req.ContentType == nil {
		req.ContentType = make(map[string]string)
	}
	if req.Params == nil {
		req.Params = make(map[string]string)
	}

	newHttpReq, err := reqStructToHttpReq(req)
	if err != nil {
		return &domain.Response{}, err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	start := time.Now()

	httpResp, err := client.Do(newHttpReq)
	if err != nil {
		return &domain.Response{}, err
	}
	defer httpResp.Body.Close()

	responseTime := time.Since(start)

	newResp := &domain.Response{}

	err = newResp.BuildResponse(httpResp)
	if err != nil {
		return &domain.Response{}, err
	}

	newResp.ResponseTime = responseTime

	return newResp, nil
}

// TODO: change database request to name and json, update save request
func (s *httpClientService) SaveRequest(req *domain.Request) error {
	// params := database.CreateRequestParams{
	// 	ID:          generateID(),
	// 	Name:        req.URL,
	// 	Method:      req.Method,
	// 	Url:         req.URL,
	// 	ContentType: toNullString(req.ContentType),
	// 	Body:        toNullString(req.Body),
	// 	Params:      toNullString(mapToString(req.Params)),
	// 	Headers:     toNullString(mapToString(req.Headers)),
	// 	Auth:        sql.NullString{}, // No auth in domain.Request yet
	// }
	//
	// _, err := s.requestRepo.Queries.CreateRequest(context.Background(), params)
	return nil
}

func (s *httpClientService) GetSavedRequests() error {
	return nil
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(parts, " ")
}

func addParams(params map[string]string, url string) string {
	if len(params) == 0 {
		return url
	}
	var paramsSlice []string
	if len(params) != 0 {
		for k, v := range params {
			paramsSlice = append(paramsSlice, k+"="+v)
		}
	}

	formattedParams := strings.Join(paramsSlice, "&")
	return url + "?" + formattedParams

}

func reqStructToHttpReq(req *domain.Request) (*http.Request, error) {

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}
	urlWithParams := addParams(req.Params, req.URL)

	httpRequest, err := http.NewRequestWithContext(context.Background(), req.Method, urlWithParams, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, val := range req.Headers {
		httpRequest.Header.Add(key, val)
	}
	if contentType, exists := req.ContentType["Content-Type"]; exists {
		httpRequest.Header.Add("Content-Type", contentType)
	} else {
		httpRequest.Header.Add("Content-Type", "text/plain; charset=utf-8")
	}

	return httpRequest, nil
}
