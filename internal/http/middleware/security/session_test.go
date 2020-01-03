//nolint:funlen // There is no sense in limiting the size of a test function
package security_test

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionHandler(t *testing.T) {
	tests := []struct {
		name           string
		preinitialized bool
	}{
		{
			name: "Base",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manager := scs.New()
			handler := security.SessionHandler{
				SessionManager: manager,
				ErrHandler:     httperror.HTTPErrorHandler{},
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			}

			recorder := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost:8080", strings.NewReader(""))

			handler.ServeHTTP(recorder, r)

			cookieHeader := recorder.Header().Get("Set-Cookie")
			_, err := asCookie(cookieHeader)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func asCookie(cookieHeader string) (http.Cookie, error) {
	splitted := strings.Split(cookieHeader, ";")
	sessionToken := strings.Trim(splitted[0], " ")
	sessionAndToken := strings.Split(sessionToken, "=")
	if sessionAndToken[0] != "session" {
		return http.Cookie{}, fmt.Errorf("header didn't parse as expected, got '%s' instead of '%s'", sessionAndToken[0], "session")
	}
	return http.Cookie{
		Name:  sessionAndToken[0],
		Value: sessionAndToken[1],
	}, nil
}

func TestAuthenticatedUserSessionUpdater_RenewSessionOnAuthenticatedUser(t *testing.T) {
	sessionManager := scs.New()
	r := httptest.NewRequest("GET", "http://localhost:8080", strings.NewReader(""))
	ctx, err := sessionManager.Load(r.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	r = r.WithContext(ctx)

	const expectedUsername = "testuser"
	sessionUpdater := security.AuthenticatedUserSessionUpdater{SessionManager: sessionManager}
	if err := sessionUpdater.RenewSessionOnAuthenticatedUser(r, expectedUsername); nil != err {
		t.Fatal(err)
	}

	user := sessionManager.GetString(r.Context(), "user")

	if expectedUsername != user {
		t.Errorf("Wrong username {expected:%s;got:%s}", expectedUsername, user)
	}
}
