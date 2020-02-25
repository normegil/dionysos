package security

import "errors"

func IsInvalidPassword(err error) bool {
	return errors.As(err, &invalidPasswordError{})
}
