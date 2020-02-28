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
	rt.Get("/sign-out", c.signOut)
	return rt
}

func (c AuthController) signOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
