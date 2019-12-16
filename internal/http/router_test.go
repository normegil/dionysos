package http_test

import (
	"fmt"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	expected := http.StatusNoContent

	routes := make(map[string]http.Handler)
	routes["/api"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expected)
	})
	srv := httptest.NewServer(internalHTTP.NewRouter(routes))
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api", srv.URL))
	if nil != err {
		t.Fatal(fmt.Sprintf(""))
	}
	if expected != resp.StatusCode {
		t.Errorf("Wrong status code {Expected:%d;Got:%d}", expected, resp.StatusCode)
	}
}
