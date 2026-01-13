package domain

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	StatusCode    int
	ContentType   string
	ContentLenght int
	Body          string
	ResponseTime  time.Duration
}

func (resp *Response) BuildResponse(httpR *http.Response) error {
	resp.StatusCode = httpR.StatusCode
	resp.ContentType = httpR.Header.Get("Content-Type")
	resp.ContentLenght = int(httpR.ContentLength)

	if httpR.Body != nil {
		bodyBytes, err := io.ReadAll(httpR.Body)
		if err != nil {
			resp.Body = fmt.Sprintf("Error reading body: %v", err)
			return err
		}
		resp.Body = string(bodyBytes)
		return nil
	}

	resp.Body = ""

	return nil
}
