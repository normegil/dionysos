package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func NewRouter(routes map[string]http.Handler) http.Handler {
	rt := chi.NewRouter()
	for pattern, handler := range routes {
		rt.Mount(pattern, handler)
	}
	return rt
}
