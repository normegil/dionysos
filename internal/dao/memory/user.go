package memory

import (
	"github.com/normegil/dionysos/internal/dao"
	"github.com/normegil/dionysos/internal/security"
)

type UserDAO struct {
	Users []*security.User
}

func (d UserDAO) Load(username string) (*security.User, error) {
	for _, user := range d.Users {
		if user.Name == username {
			return user, nil
		}
	}
	return nil, dao.NotFoundError{ID: username}
}
