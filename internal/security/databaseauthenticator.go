package security

import (
	"fmt"
)

type DatabaseAuthentication struct {
	DAO UserDAO
}

func (a DatabaseAuthentication) Authenticate(username string, password string) (*User, error) {
	user, err := a.DAO.Load(username)
	if err != nil {
		if a.DAO.IsUserNotExistError(err) {
			return nil, userNotExistError{
				Username: username,
				Original: err,
			}
		}
		return nil, fmt.Errorf("loading user '%s': %w", username, err)
	}
	err = user.ValidatePassword(password)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return user, nil
}

type UserDAO interface {
	Load(username string) (*User, error)
	IsUserNotExistError(err error) bool
}
