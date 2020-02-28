package security

import (
	"errors"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
)

const KeyUser string = "authenticated-user"

type Authenticator interface {
	Authenticate(username string, password string) (*security.User, error)
}

type RequestAuthenticator struct {
	Authenticator   Authenticator
	OnAuthenticated func(r *http.Request, username string) error
}

func (a RequestAuthenticator) Authenticate(r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if ok {
		user, err := a.Authenticator.Authenticate(username, password)
		if err != nil && !security.IsInvalidPassword(err) && !security.IsUserNotExistError(err) {
			return fmt.Errorf("error during authentication: %w", err)
		}
		if nil != user {
			if nil != a.OnAuthenticated {
				if err := a.OnAuthenticated(r, username); nil != err {
					return fmt.Errorf("authenticater user event error: %w", err)
				}
			}
		} else {
			return httperror.HTTPError{
				Code:   40100,
				Status: http.StatusUnauthorized,
				Err:    errors.New("authentication failed: wrong user and/or password"),
			}
		}
	}
	return nil
}

type AuthenticationHandler struct {
	ErrorHandler         httperror.HTTPErrorHandler
	RequestAuthenticator RequestAuthenticator
	Handler              http.Handler
}

func (a AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := a.RequestAuthenticator.Authenticate(r)
	if err != nil {
		a.ErrorHandler.Handle(w, err)
		return
	}
	a.Handler.ServeHTTP(w, r)
}
