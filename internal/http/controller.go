package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

type MultiController struct {
	Routes     map[string]http.Handler
	OnRegister func(rt *chi.Mux)
}

func (c MultiController) Route() http.Handler {
	rt := chi.NewRouter()
	for pattern, handler := range c.Routes {
		rt.Mount(pattern, handler)
	}

	if nil != c.OnRegister {
		c.OnRegister(rt)
	}
	return rt
}
