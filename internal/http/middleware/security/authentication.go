package security

import (
	"context"
	"errors"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
)

const KeyUser string = "authenticated-user"
const AnonymousUser string = "anonymous"

type Authenticator interface {
	Authenticate(username string, password string) (bool, error)
}

type AuthenticationHandler struct {
	ErrorHandler    httperror.HTTPErrorHandler
	Authenticator   Authenticator
	OnAuthenticated func(username string) error
	Handler         http.Handler
}

func (a AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		r = r.WithContext(context.WithValue(r.Context(), KeyUser, AnonymousUser))
	} else {
		authenticated, err := a.Authenticator.Authenticate(username, password)
		if err != nil && !security.IsInvalidPassword(err) {
			a.ErrorHandler.Handle(w, fmt.Errorf("error during authentication: %w", err))
			return
		}
		if authenticated {
			r = r.WithContext(context.WithValue(r.Context(), KeyUser, username))
			if nil != a.OnAuthenticated {
				if err := a.OnAuthenticated(username); nil != err {
					a.ErrorHandler.Handle(w, fmt.Errorf("authenticater user event error: %w", err))
					return
				}
			}
		} else {
			a.ErrorHandler.Handle(w, httperror.HTTPError{
				Code:   40100,
				Status: http.StatusUnauthorized,
				Err:    errors.New("authentication failed: wrong user and/or password"),
			})
			return
		}
	}

	a.Handler.ServeHTTP(w, r)
}
