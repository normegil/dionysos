package error

import "time"

type HTTPError struct {
	Code   int       `json:"code"`
	Status int       `json:"status"`
	Error  string    `json:"error"`
	Time   time.Time `json:"time"`
}
