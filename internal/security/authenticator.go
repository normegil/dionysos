package security

import (
	"fmt"
	"github.com/normegil/dionysos/internal/dao"
)

type Authenticator struct {
	DAO UserDAO
}

func (a Authenticator) Authenticate(username string, password string) (*User, error) {
	user, err := a.DAO.Load(username)
	if err != nil {
		if dao.IsNotFoundError(err) {
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
}
