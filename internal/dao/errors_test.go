package dao_test

import (
	"errors"
	"fmt"
	"github.com/normegil/dionysos/internal/dao"
	"testing"
)

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name                string
		err                 error
		expectNotFoundError bool
	}{
		{
			name: "Is a not found error",
			err: dao.NotFoundError{
				ID:       "test",
				Original: nil,
			},
			expectNotFoundError: true,
		},
		{
			name: "Is a not found error - In other errors",
			err: fmt.Errorf("containing: %w", dao.NotFoundError{
				ID:       "test",
				Original: nil,
			}),
			expectNotFoundError: true,
		},
		{
			name:                "Is a random error",
			err:                 errors.New(""),
			expectNotFoundError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isANotFoundError := dao.IsNotFoundError(test.err)
			if isANotFoundError != test.expectNotFoundError {
				t.Errorf("Results is not expected {expected:%t;got:%t}", test.expectNotFoundError, isANotFoundError)
			}
		})
	}
}
