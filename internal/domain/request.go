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
		req.URL = "http://localhost:" + cfg.DefaultPort
		return nil
	}

	if strings.HasPrefix(url, ":") {
		req.URL = "http://localhost" + url
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
		parsedHeader := strings.SplitN(h, ":", 2)
		if len(parsedHeader) == 2 {
			req.Headers[parsedHeader[0]] = parsedHeader[1]
		}
	}
	return nil
}

func (req *Request) ParseBodyType(bodyTypeStr string) error {
	req.ContentType = bodyTypeStr
	return nil
}

func (req *Request) ParseBody(body string) error {
	//TODO: add json functionality
	req.Body = body
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
		parsedParams := strings.SplitN(p, ":", 2)
		if len(parsedParams) == 2 {
			req.Params[parsedParams[0]] = parsedParams[1]
		}
	}
	return nil
}

func (req *Request) BuildRequest(method, url, headers, params, bodyType, body string, cfg *config.Config) error {
	err := req.ParseMethod(method)
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

	err = req.ParseBody(body)
	if err != nil {
		return err
	}

	return nil
}
