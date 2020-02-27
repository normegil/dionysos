package authorization

import "github.com/casbin/casbin/model"

func Model() model.Model {
	m := model.Model(make(map[string]model.AssertionMap))
	m.AddDef("r", "r", "role, resource, action")
	m.AddDef("p", "p", "role, resource, action")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.role == root || (r.role == p.role) && (r.resource == p.resource) && (r.action == p.action)")
	return m
}
