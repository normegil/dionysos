package security

type RoleRights struct {
	Role   Role              `json:"role"`
	Rights []*ResourceRights `json:"rights"`
}

func (r RoleRights) Equal(value RoleRights) bool {
	if r.Role != value.Role {
		return false
	}
	if len(r.Rights) != len(value.Rights) {
		return false
	}
	for i, rights := range r.Rights {
		if !rights.Equal(*value.Rights[i]) {
			return false
		}
	}
	return true
}

type ResourceRights struct {
	Name           string   `json:"name"`
	AllowedActions []string `json:"allowedActions"`
}

func (r ResourceRights) Equal(value ResourceRights) bool {
	if r.Name != value.Name {
		return false
	}
	if len(r.AllowedActions) != len(value.AllowedActions) {
		return false
	}
	for i, action := range r.AllowedActions {
		if action != value.AllowedActions[i] {
			return false
		}
	}
	return true
}
