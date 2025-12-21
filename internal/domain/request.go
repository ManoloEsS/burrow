package domain

import (
	"errors"
	"strings"

	"github.com/ManoloEsS/burrow/internal/config"
)

type Request struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	ContentType string            `json:"content-type,omitempty"`
	Body        string            `json:"body,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Auth        map[string]string `json:"auth,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

func (req *Request) ParseMethod(method string) error {
	if method == "" {
		return errors.New("method required for http request")
	}
	correctMethod := strings.ToUpper(method)
	req.Method = correctMethod
	return nil
}

func (req *Request) ParseUrl(cfg *config.Config, url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		req.URL = url
		return nil
	}
	if url == cfg.DefaultPort || url == "" {
		req.URL = "http://localhost" + cfg.DefaultPort
		return nil
	}

	req.URL = "https://" + url
	return nil
}

func (req *Request) ParseHeaders(headersStr string) error {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	if headersStr == "" {
		return nil
	}

	headers := strings.Fields(headersStr)
	for _, h := range headers {
		parsedHeader := strings.Split(h, ":")
		req.Headers[parsedHeader[0]] = parsedHeader[1]
	}
	return nil
}

func (req *Request) ParseBody(body string) error {
	req.Body = body
	return nil
}

func (req *Request) ParseAuth(auth string) error {
	if req.Auth == nil {
		req.Auth = make(map[string]string)
	}

	if auth == "" {
		return nil
	}

	req.Auth[auth] = "placeholder"
	return nil

}

func (req *Request) ParseParams(paramsStr string) error {
	if req.Params == nil {
		req.Params = make(map[string]string)
	}

	if paramsStr == "" {
		return nil
	}

	params := strings.Fields(paramsStr)
	for _, p := range params {
		parsedParams := strings.Split(p, ":")
		req.Params[parsedParams[0]] = parsedParams[1]
	}
	return nil
}
