package security

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/model"
	"strings"
)

type Policies []Policy

func (p Policies) ToRoleRights() []*RoleRights {
	rights := make([]*RoleRights, 0)
	for _, policy := range p {
		var foundRoleRight *RoleRights
		for _, roleRight := range rights {
			if policy.Role == roleRight.Role {
				foundRoleRight = roleRight
			}
		}
		if nil == foundRoleRight {
			foundRoleRight = &RoleRights{
				Role:   policy.Role,
				Rights: make([]*ResourceRights, 0),
			}
			rights = append(rights, foundRoleRight)
		}

		var foundRessourceRight *ResourceRights
		for _, ressourceRight := range foundRoleRight.Rights {
			if string(policy.Resource) == ressourceRight.Name {
				foundRessourceRight = ressourceRight
			}
		}
		if nil == foundRessourceRight {
			foundRessourceRight = &ResourceRights{
				Name:           string(policy.Resource),
				AllowedActions: make([]string, 0),
			}
			foundRoleRight.Rights = append(foundRoleRight.Rights, foundRessourceRight)
		}

		foundAction := false
		for _, action := range foundRessourceRight.AllowedActions {
			if string(policy.Action) == action {
				foundAction = true
			}
		}
		if !foundAction {
			foundRessourceRight.AllowedActions = append(foundRessourceRight.AllowedActions, string(policy.Action))
		}
	}

	return rights
}

type Policy struct {
	Role     Role
	Resource model.Resource
	Action   model.Action
}

type CasbinRule struct {
	ID    uuid.UUID
	Type  string
	Value string
}

func (r CasbinRule) String() string {
	const prefixLine = ", "
	var sb strings.Builder

	sb.WriteString(r.Type)
	if len(r.Value) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(r.Value)
	}
	return sb.String()
}

func (r CasbinRule) ToRule(dao RoleDAO) (*Policy, error) {
	if r.Type != "p" {
		return nil, notPolicyError{Type: r.Type}
	}

	splittedValue := strings.Split(r.Value, ", ")
	if len(splittedValue) != 3 {
		return nil, fmt.Errorf("policy has wrong format: %+v", r)
	}

	policyRoleName := splittedValue[0]
	if policyRoleName == RoleNilPolicyReference() {
		policyRoleName = RoleNil().Name
	}

	role, err := dao.LoadByName(policyRoleName)
	if err != nil {
		return nil, fmt.Errorf("loading role '%s': %w", policyRoleName, err)
	}

	return &Policy{
		Role:     *role,
		Resource: model.Resource(splittedValue[1]),
		Action:   model.Action(splittedValue[2]),
	}, nil
}

type notPolicyError struct {
	Type string
}

func (e notPolicyError) Error() string {
	return fmt.Errorf("rule is not a policy '%s'", e.Type).Error()
}

func IsANotPolicyError(err error) bool {
	return errors.As(err, &notPolicyError{})
}
