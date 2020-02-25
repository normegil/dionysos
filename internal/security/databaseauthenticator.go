package security

import (
	"fmt"
)

type DatabaseAuthentication struct {
	DAO UserDAO
}

func (a DatabaseAuthentication) Authenticate(username string, password string) (bool, error) {
	user, err := a.DAO.Load(username)
	if err != nil {
		return false, fmt.Errorf("loading user '%s': %w", username, err)
	}
	err = user.ValidatePassword(password)
	if err != nil {
		return false, fmt.Errorf("validation failed: %w", err)
	}
	return true, nil
}

type UserDAO interface {
	Load(username string) (*User, error)
}
