package domain

import "time"

type Response struct {
	ID         int
	RequestID  int
	StatusCode int
	Body       string
	Created    time.Time
}
