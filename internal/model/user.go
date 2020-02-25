package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
)

type User struct {
	ID            uuid.UUID
	Name          string
	PasswordHash  []byte
	HashAlgorithm security.HashAlgorithm
}

func NewUser(name, password string) (*User, error) {
	hashAlgorithm := security.HashAlgorithmBcrypt14
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
