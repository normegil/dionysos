package security

import (
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	PasswordHash  []byte        `json:"-"`
	HashAlgorithm HashAlgorithm `json:"-"`
	Role          Role          `json:"role"`
}

func NewUser(name, password string, role Role) (*User, error) {
	hashAlgorithm := HashAlgorithmBcrypt14
	hash, err := hashAlgorithm.HashAndSalt(password)
	if err != nil {
		return nil, fmt.Errorf("hash password of new user: %w", err)
	}
	u := &User{
		ID:            uuid.Nil,
		Name:          name,
		PasswordHash:  hash,
		HashAlgorithm: hashAlgorithm,
		Role:          role,
	}
	err = u.ValidatePassword(password)
	if err != nil {
		return nil, fmt.Errorf("validating generated hash '%s': %w", hash, err)
	}
	return u, nil
}

func (u User) ValidatePassword(password string) error {
	return u.HashAlgorithm.Validate(u.PasswordHash, password)
}

var UserAnonymous = User{
	ID:            uuid.Nil,
	Name:          "anonymous",
	PasswordHash:  []byte(""),
	HashAlgorithm: HashAlgorithmBcrypt14,
	Role:          RoleNil,
}
