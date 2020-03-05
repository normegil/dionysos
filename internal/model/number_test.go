package model_test

import (
	"github.com/normegil/dionysos/internal/model"
	"strconv"
	"testing"
)

func TestNewNatural(t *testing.T) {
	tests := []struct {
		tested      int
		expectError bool
	}{
		{
			tested:      1,
			expectError: false,
		},
		{
			tested:      0,
			expectError: false,
		},
		{
			tested:      -1,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(strconv.Itoa(test.tested), func(t *testing.T) {
			_, err := model.NewNatural(test.tested)
			if nil != err && !test.expectError {
				t.Errorf("didn't expect an error for '%d': %w", test.tested, err)
			}
		})
	}
}
