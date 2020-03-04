package security

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func RoleNil() Role {
	var role Role
	return role
}

func RoleNilPolicyReference() string {
	return "none"
}

type RoleDAO interface {
	LoadByName(string) (*Role, error)
}

type NilRoleDAO struct {
	RoleDAO RoleDAO
}

func (d NilRoleDAO) LoadByName(name string) (*Role, error) {
	roleNil := RoleNil()
	if roleNil.Name == name {
		return &roleNil, nil
	}
	return d.RoleDAO.LoadByName(name)
}
