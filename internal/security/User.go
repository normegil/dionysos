package security

import "github.com/google/uuid"

type DBUser struct {
	ID                uuid.UUID
	Name              string
	PasswordHash      string
	PasswordSalt      string
	PasswordAlgorithm string
	Roles             []Role
}
