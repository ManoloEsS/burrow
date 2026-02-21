package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
)

type httpClientService struct {
	requestRepo *database.Database
}

func NewHttpClientService(requestRepo *database.Database) HttpClientService {
	return &httpClientService{
		requestRepo: requestRepo,
	}
}

func (s *httpClientService) DeleteRequest(reqName string) error {
	err := s.requestRepo.Queries.DeleteRequest(context.Background(), reqName)
	if err != nil {
		log.Printf("could not delete request from database: %v", err)
		return fmt.Errorf("could not delete request from database: %v", err)
	}
	return nil
}

func (s *httpClientService) SendRequest(req *domain.Request) (*domain.Response, error) {
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
	defer func() { _ = httpResp.Body.Close() }()

	responseTime := time.Since(start)

	newResp := &domain.Response{}

	err = newResp.BuildResponse(httpResp)
	if err != nil {
		return &domain.Response{}, err
	}

	newResp.ResponseTime = responseTime

	return newResp, nil
}

func (s *httpClientService) SaveRequest(req *domain.Request) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request: %v", err)
	}

	requestParams := database.CreateRequestParams{
		Name:        req.Name,
		RequestJson: jsonData,
	}
	_, err = s.requestRepo.Queries.CreateRequest(context.Background(), requestParams)
	if err != nil {
		return fmt.Errorf("could not save request: %v", err)
	}

	return nil
}

func (s *httpClientService) GetSavedRequests() ([]*domain.Request, error) {
	reqsJSON, err := s.requestRepo.Queries.ListRequests(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve saved requests from database: %w", err)
	}

	var reqs []*domain.Request

	for _, r := range reqsJSON {
		request, err := requestJSONToStruct(r.RequestJson)
		if err != nil {
			log.Printf("could not parse request %s: %v", r.Name, err)
			continue
		}
		reqs = append(reqs, request)
	}
	return reqs, nil
}

func requestJSONToStruct(jsonData interface{}) (*domain.Request, error) {
	jsonByte, ok := jsonData.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", jsonData)
	}

	var req domain.Request
	err := json.Unmarshal(jsonByte, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &req, nil
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

	if req.Body == "" {
		httpRequest.Header.Add("Content-Type", "none/none")
	} else {
		httpRequest.Header.Add("Content-Type", req.ContentType["Content-Type"])
	}

	return httpRequest, nil
}
