package security

type Test_MemoryRoleDAO struct {
	Roles []*Role
}

func (d Test_MemoryRoleDAO) LoadByName(name string) (*Role, error) {
	for _, role := range d.Roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}
