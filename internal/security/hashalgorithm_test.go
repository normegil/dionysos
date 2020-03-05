package security_test

import (
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestHashAlgorithms_FindByID(t *testing.T) {
	tests := []struct {
		name          string
		algorithms    security.HashAlgorithms
		searched      uuid.UUID
		expectedFound bool
	}{
		{
			name: "Found",
			algorithms: security.HashAlgorithms([]security.HashAlgorithm{
				security.Bcrypt{
					Identifier: uuid.MustParse("e481748a-5e0f-47ac-8a29-2754168f36a5"),
					Cost:       0,
				},
			}),
			searched:      uuid.MustParse("e481748a-5e0f-47ac-8a29-2754168f36a5"),
			expectedFound: true,
		},
		{
			name: "Not found",
			algorithms: security.HashAlgorithms([]security.HashAlgorithm{
				security.Bcrypt{
					Identifier: uuid.MustParse("e481748a-5e0f-47ac-8a29-2754168f36a5"),
					Cost:       0,
				},
			}),
			searched:      uuid.Nil,
			expectedFound: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			algorithm := test.algorithms.FindByID(test.searched)
			if !test.expectedFound && algorithm != nil {
				t.Fatalf("should not found '%s' but still found it", test.searched)
			} else if test.expectedFound && algorithm == nil {
				t.Fatalf("should found '%s' but didn't find it", test.searched)
			}
		})
	}

}
