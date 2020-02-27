package authorization_test

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/normegil/dionysos/internal/security/authorization"
	"testing"
)

func TestCasbinAuthorizer_IsAuthorized(t *testing.T) {
	policyPath := "testdata/testpolicy.csv"

	tests := []struct {
		Name       string
		HTTPMethod string
		PolicyPath string
		Expected   bool
	}{
		{
			Name:       "Allowed",
			HTTPMethod: "GET",
			PolicyPath: policyPath,
			Expected:   true,
		},
		{
			Name:       "Unallowed",
			HTTPMethod: "POST",
			PolicyPath: policyPath,
			Expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			adapter := fileadapter.NewAdapter(test.PolicyPath)
			enforcer, err := casbin.NewEnforcerSafe(RESTModel(), adapter)
			if err != nil {
				t.Fatal(fmt.Errorf("error when creating casbin enforcer: %w", err))
			}
			authorizer := authorization.CasbinAuthorizer{
				Enforcer: enforcer,
			}
			authorized, err := authorizer.IsAuthorized("testuser", test.HTTPMethod, "/testuser/1")
			if err != nil {
				t.Fatal(fmt.Errorf("error when checking for authorization: %w", err))
			}
			if authorized != test.Expected {
				t.Errorf("Authorized value (%t) is not expected (%t)", authorized, test.Expected)
			}
		})
	}
}

func RESTModel() model.Model {
	m := model.Model(make(map[string]model.AssertionMap))
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")
	return m
}
