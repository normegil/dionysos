package security_test

import (
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestNilRoleDAO_LoadByName(t *testing.T) {
	tests := []struct {
		name           string
		availableRoles []*security.Role
		searched       string
		isNilRole      bool
	}{
		{
			name:           "Nil Role",
			availableRoles: []*security.Role{},
			searched:       "",
			isNilRole:      true,
		},
		{
			name: "Classic Role",
			availableRoles: []*security.Role{
				{
					ID:   uuid.New(),
					Name: "testRole",
				},
			},
			searched:  "testRole",
			isNilRole: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dao := security.NilRoleDAO{RoleDAO: security.Test_MemoryRoleDAO{Roles: test.availableRoles}}
			role, err := dao.LoadByName(test.searched)
			if err != nil {
				t.Fatal(err)
			}

			roleNil := security.RoleNil()
			if (test.isNilRole && roleNil != *role) || (!test.isNilRole && roleNil == *role) {
				t.Errorf("Not expected result {isNill:%t;got:%+v}", test.isNilRole, role)
			}
		})
	}
}
