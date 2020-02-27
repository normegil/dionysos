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
	authorized, err := c.Enforcer.EnforceSafe(role.Name, string(resource), string(action))
	if err != nil {
		return false, fmt.Errorf("loading authorizations for %s, %s %s: %w", role.Name, resource, action, err)
	}
	return authorized, nil
}
