package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/dao"
	"github.com/normegil/dionysos/internal/security"
	"strings"
)

type UserDAO struct {
	Querier Querier
}

func (d UserDAO) Load(username string) (*security.User, error) {
	row := d.Querier.QueryRow(`SELECT id, name, hash, algorithmID, roleID FROM "user" WHERE name=$1`, username)
	var id uuid.UUID
	var name string
	var hash []byte
	var algorithmID uuid.UUID
	var roleID uuid.UUID
	if err := row.Scan(&id, &name, &hash, &algorithmID, &roleID); nil != err {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, dao.NotFoundError{
				ID:       username,
				Original: err,
			}
		}
		return nil, fmt.Errorf("loading user '%s': %w", username, err)
	}
	algorithm := security.AllHashAlgorithms().FindByID(algorithmID)
	if nil == algorithm {
		return nil, fmt.Errorf("loading user '%s': algorithm not found for id '%s'", username, algorithmID)
	}

	role, err := RoleDAO(d).LoadByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("loading user '%s': %w", username, err)
	}
	return &security.User{
		ID:            id,
		Name:          name,
		PasswordHash:  hash,
		HashAlgorithm: algorithm,
		Role:          *role,
	}, nil
}

func (d UserDAO) Insert(user security.User) error {
	if _, err := d.Querier.Exec(`INSERT INTO "user" (id, name, hash, algorithmID, roleID) VALUES (gen_random_uuid(), $1, $2::bytea, $3, $4);`, user.Name, user.PasswordHash, user.HashAlgorithm.ID(), user.Role.ID); err != nil {
		return fmt.Errorf("inserting %s: %w", user.Name, err)
	}
	return nil
}
