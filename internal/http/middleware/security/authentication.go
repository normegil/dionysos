package security

import (
	"context"
	"errors"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"net/http"
)

const KeyUser string = "authenticated-user"

type Authenticator interface {
	Authenticate(username string, password string) bool
}

type AuthenticationHandler struct {
	ErrorHandler  httperror.HTTPErrorHandler
	Authenticator Authenticator
	Handler       http.Handler
}

func (a AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		a.ErrorHandler.Handle(w, httperror.HTTPError{
			Code:   40100,
			Status: http.StatusUnauthorized,
			Err:    errors.New("authentication failed: header could not be parsed"),
		})
		return
	}

	if a.Authenticator.Authenticate(username, password) {
		r = r.WithContext(context.WithValue(r.Context(), KeyUser, username))
	}

	a.Handler.ServeHTTP(w, r)
}
