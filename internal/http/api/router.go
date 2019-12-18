package api

import (
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/http/api/stock"
	"net/http"
)

func NewRouter() http.Handler {
	rt := chi.NewRouter()
	rt.Mount("/items", stock.NewController())
	rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/api", http.FileServer(pkger.Dir("/api"))).ServeHTTP(w, r)
	})
	return rt
}
