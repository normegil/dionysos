package security

import (
	"context"
	"net/http"
)

type AnonymousUserSetter struct {
	Handler http.Handler
}

func (a AnonymousUserSetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), KeyUser, AnonymousUser))
	a.Handler.ServeHTTP(w, r)
}
