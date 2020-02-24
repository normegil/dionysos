package security

import (
	"fmt"
	"github.com/casbin/casbin"
)

type CasbinAuthorizer struct {
	Enforcer *casbin.Enforcer
}

func (c CasbinAuthorizer) IsAuthorized(username string, method string, path string) (bool, error) {
	authorized, err := c.Enforcer.EnforceSafe(username, path, method)
	if err != nil {
		return false, fmt.Errorf("loading authorizations for %s, %s %s: %w", username, method, path, err)
	}
	return authorized, nil
}
