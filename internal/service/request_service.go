package service

import (
	"fmt"
	"strings"

	"github.com/ManoloEsS/burrow/internal/config"
	"github.com/ManoloEsS/burrow/internal/database"
	"github.com/ManoloEsS/burrow/internal/domain"
)

type requestService struct {
	requestRepo     *database.Database
	config          *config.Config
	updateCallback  func(*Response)
	currentResponse *Response // In-memory storage
}

// Creates new request service module for service layer
func NewRequestService(requestRepo *database.Database, config *config.Config) RequestService {
	return &requestService{
		requestRepo: requestRepo,
		config:      config,
	}
}

// Sends http request and returns response
func (s *requestService) SendRequest(req *domain.Request) (*Response, error) {

	return &Response{}, nil
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

// Gets saved requests from database
func (s *requestService) GetSavedRequests() error {
	return nil
}

// Update ui
func (s *requestService) SetUpdateCallback(callback func(*Response)) {
	s.updateCallback = callback
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

func (s *requestService) BuildRequest(method, url, headers, params, bodyType, body string) (*domain.Request, error) {
	newRequest := domain.Request{}

	err := newRequest.ParseMethod(method)
	if err != nil {
		return &newRequest, err
	}

	err = newRequest.ParseUrl(s.config, url)
	if err != nil {
		return &newRequest, err
	}

	err = newRequest.ParseHeaders(headers)
	if err != nil {
		return &newRequest, err
	}

	err = newRequest.ParseParams(params)
	if err != nil {
		return &newRequest, err
	}

	err = newRequest.ParseBodyType(bodyType)
	if err != nil {
		return &newRequest, err
	}

	err = newRequest.ParseBody(body)
	if err != nil {
		return &newRequest, err
	}

	return &newRequest, nil
}
