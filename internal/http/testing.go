package http

import (
	"github.com/normegil/dionysos/internal/http/api"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"net/http"
	"net/http/httptest"
)

func TestNewServer() *httptest.Server {
	routes := make(map[string]http.Handler)
	routes["/api"] = api.Controller{ErrHandler: error2.HTTPErrorHandler{}}.Routes()
	return httptest.NewServer(NewRouter(routes))
}
