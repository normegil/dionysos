//nolint:funlen // There is no sense in limiting the size of a test function
package middleware_test

import (
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticationHandler(t *testing.T) {
	tests := []struct {
		name                   string
		headerContent          string
		authenticationRequired bool
		expected               int
	}{
		{
			name:                   "Working case",
			headerContent:          "Basic dXNlcjpwYXNz", // user:pass
			authenticationRequired: true,
			expected:               http.StatusOK,
		},
		{
			name:                   "Wrong password - Authentication not required",
			headerContent:          "Basic dXNlcjpub3QtcGFzcw==", // user:not-pass
			authenticationRequired: false,
			expected:               http.StatusOK,
		},
		{
			name:                   "Wrong password - Authentication required",
			headerContent:          "Basic dXNlcjpub3QtcGFzcw==", // user:not-pass
			authenticationRequired: true,
			expected:               http.StatusUnauthorized,
		},
		{
			name:                   "Empty header",
			headerContent:          "",
			authenticationRequired: false,
			expected:               http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := middleware.AuthenticationHandler{
				ErrorHandler: httperror.HTTPErrorHandler{},
				Authenticator: security.MemoryAuthenticator{
					Username: "user",
					Password: "pass",
				},
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !test.authenticationRequired {
						w.WriteHeader(http.StatusOK)
						return
					}
					user := r.Context().Value(middleware.KeyUser)
					if nil == user {
						w.WriteHeader(http.StatusUnauthorized)
					}
					w.WriteHeader(http.StatusOK)
				}),
			}

			recorder := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost:8080", strings.NewReader(""))
			r.Header.Add("Authorization", test.headerContent)

			handler.ServeHTTP(recorder, r)

			resp := recorder.Result()
			if test.expected != resp.StatusCode {
				t.Errorf("Wrong status code {expected:%d;got:%d}", test.expected, resp.StatusCode)
			}
		})
	}
}
