package security

import (
	"fmt"
	"github.com/normegil/dionysos/internal/model"
)

type DatabaseAuthentication struct {
	DAO UserDAO
}

func (a DatabaseAuthentication) Authenticate(username string, password string) (bool, error) {
	user, err := a.DAO.Load(username)
	if err != nil {
		return false, fmt.Errorf("loading user '%s': %w", username, err)
	}
	return user.ValidatePassword(password), nil
}

type UserDAO interface {
	Load(username string) (model.User, error)
}
