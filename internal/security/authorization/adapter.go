package authorization

import (
	"errors"
	"fmt"
	casbinmodel "github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/model"
	"strings"
)

type Adapter struct {
	DAO PolicyDAO
}

func (a *Adapter) LoadPolicy(m casbinmodel.Model) error {
	policies, err := a.DAO.LoadAll()
	if err != nil {
		return fmt.Errorf("loading policies: %w", err)
	}

	for _, policy := range policies {
		persist.LoadPolicyLine(policy.String(), m)
	}
	return nil
}

func (a *Adapter) SavePolicy(m casbinmodel.Model) error {
	if err := a.DAO.DeleteAll(); nil != err {
		return fmt.Errorf("delete all rules: %w", err)
	}

	rules := make([]model.CasbinRule, 0)
	for _, assertionMap := range m {
		for ptype, assertion := range assertionMap {
			for _, ruleParts := range assertion.Policy {
				value := strings.Join(ruleParts, ", ")
				rules = append(rules, model.CasbinRule{
					ID:    uuid.Nil,
					Type:  ptype,
					Value: value,
				})
			}
		}
	}

	if err := a.DAO.InsertMultiple(rules); nil != err {
		return fmt.Errorf("save all rules: %w", err)
	}
	return nil
}

func (a *Adapter) AddPolicy(_ string, ptype string, rule []string) error {
	casbinRule := model.CasbinRule{
		ID:    uuid.Nil,
		Type:  ptype,
		Value: strings.Join(rule, ", "),
	}
	if err := a.DAO.Insert(casbinRule); nil != err {
		return fmt.Errorf("saving %+v: %w", rule, err)
	}
	return nil
}

func (a Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (a Adapter) RemoveFilteredPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

type PolicyDAO interface {
	LoadAll() ([]model.CasbinRule, error)
	InsertMultiple([]model.CasbinRule) error
	Insert(model.CasbinRule) error
	DeleteAll() error
}
