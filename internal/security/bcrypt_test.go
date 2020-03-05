package security_test

import (
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestBcrypt(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Simple",
			password: "simple",
		},
		{
			name:     "Trimmed but Space",
			password: "simple with space",
		},
		{
			name:     "Leading Space",
			password: " simple",
		},
		{
			name:     "Trailing Space",
			password: "simple ",
		},
		{
			name:     "Special characters",
			password: "s1[$@3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bcrypt := security.Bcrypt{
				Identifier: uuid.UUID{},
				Cost:       0, // Keep cost low for testing purpose
			}

			hash, err := bcrypt.HashAndSalt(test.password)
			if err != nil {
				t.Fatal(err)
			}

			if err = bcrypt.Validate(hash, test.password); nil != err {
				t.Fatal(err)
			}
		})
	}
}
