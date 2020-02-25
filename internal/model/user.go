package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uuid.UUID
	Name              string
	PasswordHash      string
	PasswordSalt      string
	PasswordAlgorithm HashAlgorithm
}

func NewUser(name, password string) *User {
	return &User{
		ID:                uuid.Nil,
		Name:              name,
		PasswordHash:      "",
		PasswordSalt:      "",
		PasswordAlgorithm: ,
	}
}

func (u User) ValidatePassword(password string) bool {

}
