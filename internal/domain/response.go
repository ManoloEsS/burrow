package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Status        string
	ContentType   string
	ContentLenght int64
	Body          string
	ResponseTime  time.Duration
}

func (resp *Response) BuildResponse(httpR *http.Response) error {
	resp.Status = httpR.Status
	resp.ContentType = httpR.Header.Get("Content-Type")
	resp.ContentLenght = httpR.ContentLength

	if httpR.Body != nil {
		bodyBytes, err := io.ReadAll(httpR.Body)
		if err != nil {
			resp.Body = fmt.Sprintf("Error reading body: %v", err)
			return err
		}

		if strings.HasPrefix(resp.ContentType, "application/json") {
			var prettyJson bytes.Buffer
			if err := json.Indent(&prettyJson, bodyBytes, "", " "); err != nil {
				resp.Body = fmt.Sprintf("Error prettyfying JSON body: %v", err)
				return err
			}
			resp.Body = prettyJson.String()
			return nil
		}

		resp.Body = string(bodyBytes)
		return nil
	}

	resp.Body = ""

	return nil
}
