package security_test

import (
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestUser_ValidatePassword(t *testing.T) {
	password := "user"
	user, err := security.NewUser("user", password, security.RoleNil())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		passwordToValidate string
		expected           bool
	}{
		{
			name:               "Right password",
			passwordToValidate: password,
			expected:           true,
		},
		{
			name:               "Wrong password",
			passwordToValidate: "wrong",
			expected:           false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := user.ValidatePassword(test.passwordToValidate)
			if test.expected && err != nil {
				t.Errorf("Wrong validation {expected:%t}: %w", test.expected, err)
			}
		})
	}
}
