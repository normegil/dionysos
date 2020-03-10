package test

import (
	"io"
	"net/http"
	"testing"
)

func Request(t testing.TB, method, path string, body io.Reader) *http.Request {
	const (
		baseUrl = "http://localhost"
	)
	request, err := http.NewRequest(method, baseUrl+path, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}
