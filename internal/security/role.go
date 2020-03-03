package security

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

var RoleNil = Role{
	ID:   uuid.Nil,
	Name: "none",
}

type RoleDAO interface {
	LoadByName(string) (*Role, error)
}
