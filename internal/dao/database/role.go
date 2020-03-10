package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
)

type RoleDAO struct {
	Querier Querier
}

func (d RoleDAO) LoadByID(searchedID uuid.UUID) (*security.Role, error) {
	row := d.Querier.QueryRow(`SELECT id, name FROM "role" WHERE id=$1`, searchedID)
	role, err := d.toRole(row)
	if err != nil {
		return nil, fmt.Errorf("loading role by id '%s': %w", searchedID, err)
	}
	return role, err
}

func (d RoleDAO) toRole(row *sql.Row) (*security.Role, error) {
	var id uuid.UUID
	var name string
	if err := row.Scan(&id, &name); nil != err {
		return nil, err
	}
	return &security.Role{
		ID:   id,
		Name: name,
	}, nil
}

func (d RoleDAO) LoadByName(rolename string) (*security.Role, error) {
	row := d.Querier.QueryRow(`SELECT id, name FROM "role" WHERE name=$1`, rolename)
	role, err := d.toRole(row)
	if err != nil {
		return nil, fmt.Errorf("loading role by name '%s': %w", rolename, err)
	}
	return role, err
}

func (d RoleDAO) Insert(role security.Role) error {
	_, err := d.Querier.Exec(`INSERT INTO role (id, name) VALUES (gen_random_uuid(), $1)`, role.Name)
	if err != nil {
		return fmt.Errorf("inserting %+v: %w", role, err)
	}
	return nil
}
