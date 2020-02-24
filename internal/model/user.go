package model

import "github.com/google/uuid"

type User struct {
	ID                uuid.UUID
	Name              string
	PasswordHash      string
	PasswordSalt      string
	PasswordAlgorithm string
}

func NewUser(name, password string) *User {
	return &User{
		ID:                uuid.Nil,
		Name:              name,
		PasswordHash:      "",
		PasswordSalt:      "",
		PasswordAlgorithm: "",
	}
}

func (u User) ValidatePassword(password string) bool {

}
