package api

import (
	"github.com/go-chi/chi"
	httperror "github.com/normegil/dionysos/internal/http/error"
	middlewaresecurity "github.com/normegil/dionysos/internal/http/middleware/security"
	"net/http"
)

type AuthController struct {
	ErrHandler           httperror.HTTPErrorHandler
	RequestAuthenticator middlewaresecurity.RequestAuthenticator
}

func (c AuthController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/sign-in", c.signIn)
	return rt
}

func (c AuthController) signIn(w http.ResponseWriter, _ *http.Request) {
	// Authentication handled by middlewares
	w.WriteHeader(http.StatusOK)
}
