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

type NilRoleDAO struct {
	RoleDAO RoleDAO
}

func (d NilRoleDAO) LoadByName(name string) (*Role, error) {
	if RoleNil.Name == name {
		return &RoleNil, nil
	}
	return d.RoleDAO.LoadByName(name)
}
