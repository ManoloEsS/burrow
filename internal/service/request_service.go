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

type requestService struct {
	requestRepo     *database.Database
	config          *config.Config
	updateCallback  RequestUpdateCallBack
	currentResponse *domain.Response
}

func NewRequestService(requestRepo *database.Database, config *config.Config, callback RequestUpdateCallBack) RequestService {
	return &requestService{
		requestRepo:    requestRepo,
		config:         config,
		updateCallback: callback,
	}
}

func (s *requestService) SendRequest(req *domain.Request) (*domain.Response, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}
	httpRequest, err := http.NewRequestWithContext(context.Background(), req.Method, req.URL, bodyReader)
	if err != nil {
		return &domain.Response{}, err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	start := time.Now()

	httpResp, err := client.Do(httpRequest)
	if err != nil {
		return &domain.Response{}, err
	}

	defer httpResp.Body.Close()

	responseTime := time.Since(start)

	newResp := &domain.Response{}

	err = newResp.BuildResponse(httpResp)
	if err != nil {
		// If response building fails, still set response time and call callback
		newResp.ResponseTime = responseTime
		if s.updateCallback != nil {
			s.updateCallback(newResp)
		}
		return newResp, err
	}

	newResp.ResponseTime = responseTime

	// Call callback if it exists to notify about the response
	if s.updateCallback != nil {
		s.updateCallback(newResp)
	}

	return newResp, nil
}

// TODO: change database request to name and json, update save request
func (s *requestService) SaveRequest(req *domain.Request) error {
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

func (s *requestService) GetSavedRequests() error {
	return nil
}

func (s *requestService) ResponseStringBuilder(req *domain.Request) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Mock Response for %s %s\n\n", req.Method, req.URL))
	builder.WriteString("Status: 200 OK\n\n")
	builder.WriteString("Headers:\n")
	builder.WriteString("  Content-Type: application/json\n")
	builder.WriteString("  Server: MockServer/1.0\n\n")

	if req.Body != "" {
		builder.WriteString("Echo Body:\n")
		builder.WriteString(req.Body)
	} else {
		builder.WriteString(`{"message": "This is a mock response", "method": "` + req.Method + `"}`)
	}

	return builder.String()
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
