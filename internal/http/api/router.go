package api

import (
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/http/api/stock"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"net/http"
)

type Controller struct {
	ErrHandler error2.HTTPErrorHandler
}

func (c Controller) Routes() http.Handler {
	rt := chi.NewRouter()
	rt.Mount("/items", stock.Controller{ErrHandler: c.ErrHandler}.Routes())
	rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/api", http.FileServer(pkger.Dir("/api"))).ServeHTTP(w, r)
	})
	return rt
}
