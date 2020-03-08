package authorization

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
)

type Authorizer interface {
	IsAuthorized(role security.Role, resource model.Resource, path model.Action) (bool, error)
}

type CasbinAuthorizer struct {
	Enforcer *casbin.Enforcer
}

func (c CasbinAuthorizer) IsAuthorized(role security.Role, resource model.Resource, action model.Action) (bool, error) {
	roleName := role.Name
	if security.RoleNil() == role {
		roleName = security.RoleNilPolicyReference()
	}
	authorized, err := c.Enforcer.EnforceSafe(roleName, string(resource), string(action))
	if err != nil {
		return false, fmt.Errorf("loading authorizations for '%s, %s %s': %w", roleName, resource, action, err)
	}
	return authorized, nil
}
