//nolint:funlen // There is no sense in limiting the size of a test function
package security_test

import (
	"context"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware/security"
	security2 "github.com/normegil/dionysos/internal/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSessionHandler(t *testing.T) {
	tests := []struct {
		name                   string
		user                   string
		addCookieToManager     bool
		sessionIDCookie        string
		authenticationRequired bool
		sessionIDCreated       bool
		expected               int
	}{
		{
			name:                   "Authentication with user",
			user:                   "test",
			authenticationRequired: true,
			sessionIDCreated:       true,
			expected:               http.StatusOK,
		},
		{
			name:                   "Authentication with session",
			addCookieToManager:     true,
			sessionIDCookie:        "sessID",
			authenticationRequired: true,
			sessionIDCreated:       false,
			expected:               http.StatusOK,
		},
		{
			name:                   "Authentication with user - session cookie attached",
			user:                   "test",
			addCookieToManager:     true,
			sessionIDCookie:        "sessID",
			authenticationRequired: true,
			sessionIDCreated:       true,
			expected:               http.StatusOK,
		},
		{
			name:                   "Authentication - No data",
			authenticationRequired: true,
			sessionIDCreated:       false,
			expected:               http.StatusUnauthorized,
		},
		{
			name:                   "Authentication - No data - Not required",
			authenticationRequired: false,
			sessionIDCreated:       false,
			expected:               http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionManager := security2.NewMemorySessionManager()
			handler := security.SessionHandler{
				ErrHandler:        httperror.HTTPErrorHandler{},
				SessionTimeToLive: 1 * time.Second,
				SessionManager:    sessionManager,
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !test.authenticationRequired {
						w.WriteHeader(http.StatusOK)
						return
					}
					user := r.Context().Value(security.KeyUser)
					if nil == user {
						w.WriteHeader(http.StatusUnauthorized)
					}
					w.WriteHeader(http.StatusOK)
				}),
			}
			recorder := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost:8080", strings.NewReader(""))

			if "" != test.user {
				r = r.WithContext(context.WithValue(r.Context(), security.KeyUser, test.user))
			}

			if "" != test.sessionIDCookie {
				if test.addCookieToManager {
					sessionManager.TestSetSession("test", test.sessionIDCookie)
				}
				r.AddCookie(&http.Cookie{
					Name:  security.KeySessionID,
					Value: test.sessionIDCookie,
				})
			}

			handler.ServeHTTP(recorder, r)

			resp := recorder.Result()
			checkCookieExist(t, test.sessionIDCreated, recorder)
			checkStatusCode(t, test.expected, resp)
		})
	}
}

func checkStatusCode(t *testing.T, expected int, resp *http.Response) {
	if expected != resp.StatusCode {
		t.Errorf("Wrong status code {expected:%d,got:%d}", expected, resp.StatusCode)
	}
}

func checkCookieExist(t *testing.T, shouldExist bool, recorder *httptest.ResponseRecorder) {
	if shouldExist {
		cookieID, _ := GetCookie(t, recorder)
		if cookieID != security.KeySessionID {
			t.Errorf("no cookie '%s' found", security.KeySessionID)
			return
		}
	}
}

func GetCookie(t *testing.T, recorder *httptest.ResponseRecorder) (string, string) {
	cookieHeader := recorder.Header().Get("Set-Cookie")
	if "" == cookieHeader {
		t.Errorf("no cookie in response")
		return "", ""
	}

	splittedCookieHeader := strings.Split(cookieHeader, ";")
	cookieKeyVal := strings.Trim(splittedCookieHeader[0], " ")
	keyVal := strings.Split(cookieKeyVal, "=")
	return keyVal[0], keyVal[1]
}
