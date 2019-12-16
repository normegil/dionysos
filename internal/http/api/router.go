package api

import (
	"github.com/go-chi/chi"
	"github.com/normegil/dionysos/internal/http/api/stock"
	"net/http"
)

func NewRouter() http.Handler {
	rt := chi.NewRouter()
	rt.Mount("/items", stock.NewController())
	return rt
}
