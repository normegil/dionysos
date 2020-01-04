package security

import (
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"net/http"
)

type Authorizer interface {
	IsAuthorized(username string, method string, path string) (bool, error)
}

type AuthorizationHandler struct {
	Authorizer Authorizer
	ErrHandler httperror.HTTPErrorHandler
	Handler    http.Handler
}

func (a AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(KeyUser).(string)
	authorized, err := a.Authorizer.IsAuthorized(username, r.Method, r.RequestURI)
	if err != nil {
		a.ErrHandler.Handle(w, err)
		return
	}
	if !authorized {
		a.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40300,
			Status: http.StatusForbidden,
			Err:    fmt.Errorf("access forbidden {method:%s;ressource:%s}", r.Method, r.RequestURI),
		})
	}
}
