package security

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"net/http"
	"time"
)

const keySessionUser = "user"

type SessionHandler struct {
	SessionManager       *scs.SessionManager
	RequestAuthenticator RequestAuthenticator
	ErrHandler           httperror.HTTPErrorHandler
	Handler              http.Handler
}

func (s SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var token string
	cookie, err := r.Cookie(s.SessionManager.Cookie.Name)
	if err == nil {
		token = cookie.Value
	}

	ctx, err := s.SessionManager.Load(r.Context(), token)
	if err != nil {
		s.ErrHandler.Handle(w, fmt.Errorf("could not load session: %w", err))
		return
	}

	user := s.SessionManager.Get(ctx, keySessionUser)
	if nil != user {
		ctx = context.WithValue(ctx, KeyUser, user)
	}

	sr := r.WithContext(ctx)
	if err := s.RequestAuthenticator.Authenticate(sr); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	switch s.SessionManager.Status(ctx) {
	case scs.Unmodified:
		fallthrough
	case scs.Modified:
		token, expiry, err := s.SessionManager.Commit(ctx)
		if err != nil {
			s.ErrHandler.Handle(w, fmt.Errorf("could not commit session: %w", err))
			return
		}
		s.writeSession(w, token, expiry)
	case scs.Destroyed:
		s.writeSession(w, "", time.Time{})
	}

	s.Handler.ServeHTTP(w, sr)
}

func (s SessionHandler) writeSession(w http.ResponseWriter, token string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     s.SessionManager.Cookie.Name,
		Value:    token,
		Path:     s.SessionManager.Cookie.Path,
		Domain:   s.SessionManager.Cookie.Domain,
		Secure:   s.SessionManager.Cookie.Secure,
		HttpOnly: s.SessionManager.Cookie.HttpOnly,
		SameSite: s.SessionManager.Cookie.SameSite,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if s.SessionManager.Cookie.Persist {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.
	}

	w.Header().Add("Set-Cookie", cookie.String())
	addHeaderIfMissing(w, "Cache-Control", `no-cache="Set-Cookie"`)
	addHeaderIfMissing(w, "Vary", "Cookie")
}

func addHeaderIfMissing(w http.ResponseWriter, key, value string) {
	for _, h := range w.Header()[key] {
		if h == value {
			return
		}
	}
	w.Header().Add(key, value)
}

type AuthenticatedUserSessionUpdater struct {
	SessionManager *scs.SessionManager
}

func (a AuthenticatedUserSessionUpdater) RenewSessionOnAuthenticatedUser(r *http.Request, username string) error {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		return fmt.Errorf("could not renew session token: %w", err)
	}
	a.SessionManager.Put(r.Context(), keySessionUser, username)
	return nil
}
