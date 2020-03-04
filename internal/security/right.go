package security

type RoleRights struct {
	Role   Role             `json:"role"`
	Rights []*ResourceRights `json:"rights"`
}

type ResourceRights struct {
	Name           string   `json:"name"`
	AllowedActions []string `json:"allowedActions"`
}
