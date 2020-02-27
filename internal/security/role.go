package security

type Role struct {
	Name string `json:"name"`
}

var RoleNone = Role{Name:"none"}