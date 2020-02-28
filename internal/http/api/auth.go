package api

import (
	"github.com/go-chi/chi"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"net/http"
)

type AuthController struct {
	ErrHandler     httperror.HTTPErrorHandler
	UserController UserController
}

func (c AuthController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/sign-in", c.UserController.current)
	return rt
}
