package security

import (
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/dionysos/internal/security/authorization"
	"net/http"
)

type AuthorizationHandler struct {
	Authorizer authorization.Authorizer
	ErrHandler httperror.HTTPErrorHandler
}

func (a AuthorizationHandler) Register(resource model.Resource, action model.Action, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(KeyUser).(security.User)
		authorized, err := a.Authorizer.IsAuthorized(user.Role, resource, action)
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
			return
		}
		handler.ServeHTTP(w, r)
	}
}
