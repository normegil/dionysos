package security_test

import (
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestMemoryAuthenticator_Authenticate(t *testing.T) {
	tests := []struct {
		name                string
		authenticator       security.MemoryAuthenticator
		user                string
		password            string
		expectAuthenticated bool
	}{
		{
			name: "Authenticated",
			authenticator: security.MemoryAuthenticator{
				Username: "test",
				Password: "test",
			},
			user:                "test",
			password:            "test",
			expectAuthenticated: true,
		},
		{
			name: "Wrong user",
			authenticator: security.MemoryAuthenticator{
				Username: "test",
				Password: "test",
			},
			user:                "wrong",
			password:            "test",
			expectAuthenticated: false,
		},
		{
			name: "Wrong password",
			authenticator: security.MemoryAuthenticator{
				Username: "test",
				Password: "test",
			},
			user:                "test",
			password:            "wrong",
			expectAuthenticated: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authenticated, err := test.authenticator.Authenticate(test.user, test.password)
			if err != nil {
				t.Fatal(err)
			}
			if authenticated != test.expectAuthenticated {
				t.Errorf("Authentication didn't have expected result {expected:%t;got:%t}", test.expectAuthenticated, authenticated)
			}
		})
	}
}
