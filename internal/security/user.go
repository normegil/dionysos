package security

import (
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Name          string
	PasswordHash  []byte
	HashAlgorithm HashAlgorithm
	Role          Role
}

func NewUser(name, password string) (*User, error) {
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
