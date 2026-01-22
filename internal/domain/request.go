package domain

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ManoloEsS/burrow/internal/config"
)

type Request struct {
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	ContentType map[string]string `json:"content-type,omitempty"`
	Body        string            `json:"body,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

func NewRequest() *Request {
	return &Request{
		ContentType: make(map[string]string),
		Headers:     make(map[string]string),
		Params:      make(map[string]string),
	}
}

func (req *Request) ParseMethod(method string) error {
	if method == "" {
		return errors.New("method required for http request")
	}
	correctMethod := strings.ToUpper(strings.TrimSpace(method))
	req.Method = correctMethod
	return nil
}

func (req *Request) ParseUrl(cfg *config.Config, url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		req.URL = url
		return nil
	}
	if url == cfg.App.DefaultPort || url == "" {
		req.URL = "http://localhost:" + cfg.App.DefaultPort
		return nil
	}

	if strings.HasPrefix(url, "/") {
		req.URL = "http://localhost:" + cfg.App.DefaultPort + url
		return nil
	}

	if strings.HasPrefix(url, ":") {
		req.URL = "http://localhost" + url
		return nil
	}

	if strings.HasPrefix(url, "localhost:") {
		req.URL = "http://" + url
		return nil
	}

	req.URL = "https://" + url
	return nil
}

func (req *Request) ParseHeaders(headersStr string) error {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	if headersStr != "" {
		headers := strings.Split(headersStr, ",")
		for _, h := range headers {
			trimmedHeader := strings.TrimSpace(h)
			parsedHeader := strings.SplitN(trimmedHeader, ":", 2)
			if len(parsedHeader) == 2 {
				req.Headers[parsedHeader[0]] = parsedHeader[1]
			}
		}
	}
	req.Headers["User-Agent"] = "Burrow/1.0.0(github.com/ManoloEsS/burrow)"
	return nil
}

func (req *Request) ParseBodyType(bodyTypeStr string) error {
	if req.ContentType == nil {
		req.ContentType = make(map[string]string)
	}

	if bodyTypeStr == "Text" {
		req.ContentType["Content-Type"] = "text/plain; charset=utf-8"
	}
	if bodyTypeStr == "JSON" {
		req.ContentType["Content-Type"] = "application/json"
	}

	return nil
}

func (req *Request) ParseBody(body, bodyTypeStr string) error {
	if bodyTypeStr == "JSON" {
		if json.Valid([]byte(body)) {
			req.Body = body
			return nil
		}
		return errors.New("invalid JSON string")
	}
	if bodyTypeStr == "Text" {
		req.Body = body
		return nil
	}
	return nil
}

func (req *Request) ParseParams(paramsStr string) error {
	if req.Params == nil {
		req.Params = make(map[string]string)
	}

	if paramsStr == "" {
		return nil
	}

	params := strings.Split(paramsStr, ",")
	for _, p := range params {
		trimmedParams := strings.TrimSpace(p)
		parsedParams := strings.SplitN(trimmedParams, ":", 2)
		if len(parsedParams) == 2 {
			req.Params[parsedParams[0]] = parsedParams[1]
		}
	}
	return nil
}

func (req *Request) ParseName(nameStr string) error {
	if nameStr == "" {
		return nil
	}
	req.Name = strings.ToLower(nameStr)

	return nil
}

func (req *Request) BuildRequest(name, method, url, headers, params, bodyType, body string, cfg *config.Config) error {
	err := req.ParseName(name)
	if err != nil {
		return err
	}
	err = req.ParseMethod(method)
	if err != nil {
		return err
	}

	err = req.ParseUrl(cfg, url)
	if err != nil {
		return err
	}

	err = req.ParseHeaders(headers)
	if err != nil {
		return err
	}

	err = req.ParseParams(params)
	if err != nil {
		return err
	}

	err = req.ParseBodyType(bodyType)
	if err != nil {
		return err
	}

	err = req.ParseBody(body, bodyType)
	if err != nil {
		return err
	}

	return nil
}
