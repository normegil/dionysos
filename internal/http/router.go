package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func NewRouter() http.Handler {
	return chi.NewRouter()
}
