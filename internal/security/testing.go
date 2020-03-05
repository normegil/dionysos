package security

type Test_FakeRoleDAO struct {
	Roles []*Role
}

func (d Test_FakeRoleDAO) LoadByName(name string) (*Role, error) {
	for _, role := range d.Roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}
