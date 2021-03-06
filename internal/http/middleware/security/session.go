package security

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
	"time"
)

const keySessionUser = "user"

type SessionHandler struct {
	SessionManager       *scs.SessionManager
	RequestAuthenticator RequestAuthenticator
	ErrHandler           httperror.HTTPErrorHandler
	UserDAO              security.UserDAO
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
	r = r.WithContext(ctx)

	if err := s.handleAuthenticationAction(r); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	if err := s.RequestAuthenticator.Authenticate(r); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	username := s.SessionManager.Get(ctx, keySessionUser)
	if nil != username {
		usernameStr := username.(string)
		if "" != usernameStr && security.UserAnonymous().Name != usernameStr {
			user, err := s.UserDAO.Load(usernameStr)
			if err != nil {
				s.ErrHandler.Handle(w, fmt.Errorf("could not load user '%s': %w", usernameStr, err))
				return
			}
			ctx = context.WithValue(ctx, KeyUser, *user)
		}
	}
	sr := r.WithContext(ctx)

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

func (s SessionHandler) handleAuthenticationAction(r *http.Request) error {
	authenticationAction := r.Header.Get("X-Authentication-Action")
	if authenticationAction != "" {
		userSessionUpdater := AuthenticatedUserSessionUpdater{SessionManager: s.SessionManager}
		switch authenticationAction {
		case "sign-out":
			err := userSessionUpdater.SignOut(r)
			if err != nil {
				return fmt.Errorf("couldn't sign out: %w", err)
			}
		default:
			return fmt.Errorf("unrecognized authentication action: '%s'", authenticationAction)
		}
	}
	return nil
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

func (a AuthenticatedUserSessionUpdater) SignOut(r *http.Request) error {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		return fmt.Errorf("could not renew session token: %w", err)
	}
	a.SessionManager.Put(r.Context(), keySessionUser, security.UserAnonymous().Name)
	return nil
}
