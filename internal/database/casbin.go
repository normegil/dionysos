package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/model"
	"github.com/rs/zerolog/log"
)

type CasbinDAO struct {
	DB *sql.DB
}

func (d CasbinDAO) LoadAll() ([]model.CasbinRule, error) {
	rows, err := d.DB.Query("SELECT id, type, value FROM policy")

	if err != nil {
		return nil, fmt.Errorf("select all policies: %w", err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("select all: could not close rows")
		}
	}()

	rules := make([]model.CasbinRule, 0)
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
		rules = append(rules, model.CasbinRule{
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

func (d CasbinDAO) DeleteAll() error {
	if _, err := d.DB.Exec(fmt.Sprintf(`DELETE FROM policy`)); nil != err {
		return fmt.Errorf("deleting all policies: %w", err)
	}
	return nil
}

func (d CasbinDAO) InsertMultiple(rules []model.CasbinRule) error {
	for _, rule := range rules {
		if err := d.Insert(rule); nil != err {
			return err
		}
	}
	return nil
}

func (d CasbinDAO) Insert(rule model.CasbinRule) error {
	if _, err := d.DB.Exec(`INSERT INTO policy (id, type, value) VALUES (gen_random_uuid(), $1, $2)`, rule.Type, rule.Value); nil != err {
		return fmt.Errorf("inserting %+v: %w", rule, err)
	}
	return nil
}
