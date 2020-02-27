package security

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

var RoleNone = Role{Name: "none"}
