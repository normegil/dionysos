package security

import (
	"context"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"net/http"
	"time"
)

const KeySessionID = "sessionID"

type SessionManager interface {
	CreateSession(user string) (string, error)
	GetSession(id string) (string, error)
}

type SessionHandler struct {
	ErrHandler        httperror.HTTPErrorHandler
	SessionManager    SessionManager
	SessionTimeToLive time.Duration
	Handler           http.Handler
}

func (s SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(KeyUser)
	if nil != username {
		if err := s.newSession(username.(string), w, r); nil != err {
			s.ErrHandler.Handle(w, err)
		}
		return
	}

	sessionIDCookie, err := r.Cookie(KeySessionID)
	if err != nil && err != http.ErrNoCookie {
		s.ErrHandler.Handle(w, fmt.Errorf("error retrieving cookie '%s': %w", KeySessionID, err))
	}

	if nil != sessionIDCookie {
		if err := s.loadSession(sessionIDCookie, r, w); nil != err {
			s.ErrHandler.Handle(w, err)
		}
		return
	}

	s.Handler.ServeHTTP(w, r)
}

func (s SessionHandler) loadSession(sessionIDCookie *http.Cookie, r *http.Request, w http.ResponseWriter) error {
	user, err := s.SessionManager.GetSession(sessionIDCookie.Value)
	if err != nil {
		return fmt.Errorf("could not load session: %w", err)
	}
	r = r.WithContext(context.WithValue(r.Context(), KeyUser, user))
	s.Handler.ServeHTTP(w, r)
	return nil
}

func (s SessionHandler) newSession(username string, w http.ResponseWriter, r *http.Request) error {
	session, err := s.SessionManager.CreateSession(username)
	if err != nil {
		return fmt.Errorf("could not create session: %w", err)
	}
	s.Handler.ServeHTTP(w, r)

	http.SetCookie(w, &http.Cookie{
		Name:    KeySessionID,
		Value:   session,
		Expires: time.Now().Add(s.SessionTimeToLive),
		MaxAge:  0,
	})
	return nil
}
