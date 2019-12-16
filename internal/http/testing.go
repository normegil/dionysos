package http

import (
	"github.com/normegil/dionysos/internal/http/api"
	"net/http"
	"net/http/httptest"
)

func TestNewServer() *httptest.Server {
	routes := make(map[string]http.Handler)
	routes["/api"] = api.NewRouter()
	return httptest.NewServer(NewRouter(routes))
}
