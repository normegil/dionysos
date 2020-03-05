//nolint:funlen
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

func TestPolicies_ToRoleRights(t *testing.T) {
	role1 := security.Role{
		ID:   uuid.New(),
		Name: "testrole1",
	}
	role2 := security.Role{
		ID:   uuid.New(),
		Name: "testrole2",
	}

	tests := []struct {
		name     string
		policies security.Policies
		expected []*security.RoleRights
	}{
		{
			name:     "Empty policies",
			policies: security.Policies([]security.Policy{}),
			expected: []*security.RoleRights{},
		},
		{
			name: "One role",
			policies: security.Policies([]security.Policy{
				{Role: role1, Resource: model.ResourceItem, Action: model.ActionRead},
				{Role: role1, Resource: model.ResourceItem, Action: model.ActionWrite},
				{Role: role1, Resource: model.ResourceStorage, Action: model.ActionRead},
			}),
			expected: []*security.RoleRights{
				{
					Role: role1,
					Rights: []*security.ResourceRights{
						{Name: string(model.ResourceItem), AllowedActions: []string{
							string(model.ActionRead),
							string(model.ActionWrite),
						}},
						{Name: string(model.ResourceStorage), AllowedActions: []string{string(model.ActionRead)}},
					},
				},
			},
		},
		{
			name: "Multiple roles",
			policies: security.Policies([]security.Policy{
				{Role: role1, Resource: model.ResourceItem, Action: model.ActionRead},
				{Role: role1, Resource: model.ResourceItem, Action: model.ActionWrite},
				{Role: role2, Resource: model.ResourceItem, Action: model.ActionRead},
				{Role: role1, Resource: model.ResourceStorage, Action: model.ActionRead},
				{Role: role2, Resource: model.ResourceStorage, Action: model.ActionRead},
			}),
			expected: []*security.RoleRights{
				{
					Role: role1,
					Rights: []*security.ResourceRights{
						{Name: string(model.ResourceItem), AllowedActions: []string{
							string(model.ActionRead),
							string(model.ActionWrite),
						}},
						{Name: string(model.ResourceStorage), AllowedActions: []string{string(model.ActionRead)}},
					},
				},
				{
					Role: role2,
					Rights: []*security.ResourceRights{
						{Name: string(model.ResourceItem), AllowedActions: []string{string(model.ActionRead)}},
						{Name: string(model.ResourceStorage), AllowedActions: []string{string(model.ActionRead)}},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rights := test.policies.ToRoleRights()

			if len(rights) != len(test.expected) {
				t.Fatalf("Wrong number of role rights {expected:%d;got:%d}", len(test.expected), len(rights))
			}

			for i, expected := range test.expected {
				if !expected.Equal(*rights[i]) {
					t.Fatalf("Not equal value {index:%d;expected:%+v;got:%+v}", i, expected, *rights[i])
				}
			}
		})
	}
}
