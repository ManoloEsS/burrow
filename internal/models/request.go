package models

type Request struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	ContentType string            `json:"content-type,omitempty"`
	Body        string            `json:"body,omitempty"`
	Auth        map[string]string `json:"auth,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}
