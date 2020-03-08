package authorization_test

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/dionysos/internal/security/authorization"
	"testing"
)

func TestCasbinAuthorizer_IsAuthorized(t *testing.T) {
	m := authorization.Model()
	rule := security.CasbinRule{
		ID:    uuid.UUID{},
		Type:  "p",
		Value: "test, item, read",
	}
	persist.LoadPolicyLine(rule.String(), m)

	tests := []struct {
		name             string
		role             security.Role
		resource         model.Resource
		action           model.Action
		expectAuthorized bool
	}{
		{
			name:             "Authorized",
			role:             security.Role{ID: uuid.New(), Name: "test"},
			resource:         "item",
			action:           "read",
			expectAuthorized: true,
		},
		{
			name:             "Not authorized - Wrong role",
			role:             security.Role{ID: uuid.New(), Name: "wrong"},
			resource:         "item",
			action:           "read",
			expectAuthorized: false,
		},
		{
			name:             "Not authorized - Wrong resource",
			role:             security.Role{ID: uuid.New(), Name: "test"},
			resource:         "wrong",
			action:           "read",
			expectAuthorized: false,
		},
		{
			name:             "Not authorized - Wrong action",
			role:             security.Role{ID: uuid.New(), Name: "test"},
			resource:         "item",
			action:           "wrong",
			expectAuthorized: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authorizer := authorization.CasbinAuthorizer{Enforcer: casbin.NewEnforcer(m)}
			authorized, err := authorizer.IsAuthorized(test.role, test.resource, test.action)
			if err != nil {
				t.Fatal(err)
			}
			if test.expectAuthorized != authorized {
				t.Errorf("Authorization not in expected state {expeted:%t;got:%t}", test.expectAuthorized, authorized)
			}
		})
	}
}
