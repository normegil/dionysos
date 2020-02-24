package security_test

import (
	"context"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthorizationHandler(t *testing.T) {
	tests := []struct {
		Name           string
		Response       bool
		ExpectedStatus int
	}{
		{
			Name:           "Authorized",
			Response:       true,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Unauthorized",
			Response:       false,
			ExpectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			handler := security.AuthorizationHandler{
				Authorizer: testAuthorizer{Response: test.Response},
				ErrHandler: httperror.HTTPErrorHandler{},
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			}
			recorder := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost:8080/testuser", strings.NewReader(""))
			handler.ServeHTTP(recorder, r.WithContext(context.WithValue(r.Context(), security.KeyUser, "testuser")))
			resp := recorder.Result()
			if test.ExpectedStatus != resp.StatusCode {
				t.Errorf("Wrong status code {expected:%d;got:%d}", test.ExpectedStatus, resp.StatusCode)
			}
		})
	}
}

type testAuthorizer struct {
	Response bool
}

func (a testAuthorizer) IsAuthorized(_ string, _ string, _ string) (bool, error) {
	return a.Response, nil
}
