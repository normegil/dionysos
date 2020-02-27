package security

import (
	"context"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
)

type AnonymousUserSetter struct {
	Handler http.Handler
}

func (a AnonymousUserSetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), KeyUser, security.UserAnonymous))
	a.Handler.ServeHTTP(w, r)
}
