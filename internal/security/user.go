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
}

func NewUser(name, password string) (*User, error) {
	hashAlgorithm := HashAlgorithmBcrypt14
	hash, err := hashAlgorithm.HashAndSalt(password)
	if err != nil {
		return nil, fmt.Errorf("hash password of new user")
	}
	return &User{
		ID:            uuid.Nil,
		Name:          name,
		PasswordHash:  hash,
		HashAlgorithm: hashAlgorithm,
	}, nil
}

func (u User) ValidatePassword(password string) error {
	return u.HashAlgorithm.Validate(u.PasswordHash, password)
}
