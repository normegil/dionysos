package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
	"github.com/rs/zerolog/log"
)

type CasbinDAO struct {
	Querier Querier
	RoleDAO security.RoleDAO
}

func (d CasbinDAO) LoadAll() ([]security.CasbinRule, error) {
	rows, err := d.Querier.Query("SELECT id, type, value FROM policy")

	if err != nil {
		return nil, fmt.Errorf("select all policies: %w", err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("select all: could not close rows")
		}
	}()

	rules := make([]security.CasbinRule, 0)
	for rows.Next() {
		var idStr string
		var pType string
		var value string
		if err := rows.Scan(&idStr, &pType, &value); nil != err {
			return nil, fmt.Errorf("scanning rows of policies: %w", err)
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("cannot convert policy id '%s' to UUID : %w", idStr, err)
		}
		rules = append(rules, security.CasbinRule{
			ID:    id,
			Type:  pType,
			Value: value,
		})
	}

	if err := rows.Err(); nil != err {
		return nil, fmt.Errorf("parsing rows queried by policies: %w", err)
	}

	return rules, nil
}

func (d CasbinDAO) LoadAsRolesRights() ([]*security.RoleRights, error) {
	rules, err := d.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("load all rules: %w", err)
	}

	var policies security.Policies = make([]security.Policy, 0)
	for _, casbinRule := range rules {
		rule, err := casbinRule.ToRule(d.RoleDAO)
		if err != nil && !security.IsANotPolicyError(err) {
			return nil, fmt.Errorf("load rules as role rights: %w", err)
		}

		policies = append(policies, *rule)
	}

	return policies.ToRoleRights(), nil
}

func (d CasbinDAO) LoadForUser(user security.User) ([]*security.ResourceRights, error) {
	rolesRights, err := d.LoadAsRolesRights()
	if err != nil {
		return nil, fmt.Errorf("load all rules: %w", err)
	}

	for _, roleRights := range rolesRights {
		if roleRights.Role == user.Role {
			return roleRights.Rights, nil
		}
	}
	return nil, nil
}

func (d CasbinDAO) DeleteAll() error {
	if _, err := d.Querier.Exec(fmt.Sprintf(`DELETE FROM policy`)); nil != err {
		return fmt.Errorf("deleting all policies: %w", err)
	}
	return nil
}

func (d CasbinDAO) InsertMultiple(rules []security.CasbinRule) error {
	for _, rule := range rules {
		if err := d.Insert(rule); nil != err {
			return err
		}
	}
	return nil
}

func (d CasbinDAO) Insert(rule security.CasbinRule) error {
	if _, err := d.Querier.Exec(`INSERT INTO policy (id, type, value) VALUES (gen_random_uuid(), $1, $2)`, rule.Type, rule.Value); nil != err {
		return fmt.Errorf("inserting %+v: %w", rule, err)
	}
	return nil
}
