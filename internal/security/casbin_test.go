package security_test

import (
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestCasbinRule_String(t *testing.T) {
	tests := []struct {
		name     string
		rule     security.CasbinRule
		expected string
	}{
		{
			name: "Simple",
			rule: security.CasbinRule{
				ID:    uuid.New(),
				Type:  "p",
				Value: "v1, v2",
			},
			expected: "p, v1, v2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := test.rule.String()
			if val != test.expected {
				t.Errorf("Wrong result {expected:%s;got:%s}", test.expected, val)
			}
		})
	}
}

func TestCasbinRule_ToRule(t *testing.T) {
	role := &security.Role{
		ID:   uuid.New(),
		Name: "testrole",
	}
	dao := security.Test_MemoryRoleDAO{
		Roles: []*security.Role{
			role,
		},
	}

	tests := []struct {
		name     string
		rule     security.CasbinRule
		expected security.Policy
	}{
		{
			name: "Item-Read",
			rule: security.CasbinRule{
				ID:    uuid.New(),
				Type:  "p",
				Value: role.Name + ", item, read",
			},
			expected: security.Policy{
				Role:     *role,
				Resource: model.ResourceItem,
				Action:   model.ActionRead,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := test.rule.ToRule(dao)
			if err != nil {
				t.Fatal(err)
			}
			if *val != test.expected {
				t.Errorf("Wrong result {expected:%+v;got:%+v}", test.expected, val)
			}
		})
	}

}

func TestCasbinRule_ToRule_NotAPolicy(t *testing.T) {
	role := &security.Role{
		ID:   uuid.New(),
		Name: "testrole",
	}
	dao := security.Test_MemoryRoleDAO{
		Roles: []*security.Role{
			role,
		},
	}

	rule := security.CasbinRule{
		ID:    uuid.New(),
		Type:  "g",
		Value: role.Name + ", item, read",
	}
	_, err := rule.ToRule(dao)
	if err == nil {
		t.Fatal("error expected but got none")
	} else if !security.IsANotPolicyError(err) {
		t.Errorf("Expected a not a policy error but got: %s", err.Error())
	}
}
