package security_test

import (
	"github.com/normegil/dionysos/internal/dao/memory"
	"github.com/normegil/dionysos/internal/security"
	"testing"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	user, err := security.NewUser("test", "test", security.RoleNil())
	if err != nil {
		t.Fatal(err)
	}
	dao := memory.UserDAO{Users: []*security.User{user}}

	tests := []struct {
		name                      string
		user                      string
		password                  string
		expectWrongAuthentication bool
		expectUserNotExist        bool
	}{
		{
			name:                      "User exist - Right password",
			user:                      user.Name,
			password:                  "test",
			expectUserNotExist:        false,
			expectWrongAuthentication: false,
		},
		{
			name:                      "User exist - False password",
			user:                      user.Name,
			password:                  "test",
			expectUserNotExist:        false,
			expectWrongAuthentication: true,
		},
		{
			name:                      "User not exist",
			user:                      user.Name,
			password:                  "test",
			expectUserNotExist:        true,
			expectWrongAuthentication: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authenticator := security.Authenticator{DAO: dao}
			user, err := authenticator.Authenticate(test.user, test.password)
			if test.expectUserNotExist && security.IsUserNotExistError(err) {
				return
			} else if err != nil {
				if test.expectWrongAuthentication {
					return
				}
				t.Fatal(err)
			}

			loaded, err := dao.Load(test.user)
			if err != nil {
				t.Fatal(err)
			}
			if user != loaded {
				t.Errorf("User is not expected {expected:%+v;got:%+v}", loaded, user)
			}
		})
	}
}
