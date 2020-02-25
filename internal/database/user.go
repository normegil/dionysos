package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/security"
	"strings"
)

type UserDAO struct {
	DB *sql.DB
}

func (d UserDAO) Load(username string) (*security.User, error) {
	row := d.DB.QueryRow(`SELECT id, name, hash, algorithmID FROM "user" WHERE name=$1`, username)
	var id uuid.UUID
	var name string
	var hash []byte
	var algorithmID uuid.UUID
	if err := row.Scan(&id, &name, &hash, &algorithmID); nil != err {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, userNotExistError{
				Username: username,
				Original: err,
			}
		}
		return nil, fmt.Errorf("loading user '%s': %w", username, err)
	}
	algorithm := security.AllHashAlgorithms.FindByID(algorithmID)
	if nil == algorithm {
		return nil, fmt.Errorf("loading user '%s': algorithm not found for id '%s'", username, algorithmID)
	}
	return &security.User{
		ID:            id,
		Name:          name,
		PasswordHash:  hash,
		HashAlgorithm: algorithm,
	}, nil
}

func (d UserDAO) Insert(user security.User) error {
	if _, err := d.DB.Exec(`INSERT INTO "user" (id, name, hash, algorithmID) VALUES (gen_random_uuid(), $1, $2::bytea, $3);`, user.Name, user.PasswordHash, user.HashAlgorithm.ID()); err != nil {
		return fmt.Errorf("inserting %s: %w", user.Name, err)
	}
	return nil
}

type userNotExistError struct {
	Username string
	Original error
}

func (e userNotExistError) Error() string {
	return fmt.Errorf("databse user '%s' doesn't exist: %w", e.Username, e.Original).Error()
}


func (d UserDAO) IsUserNotExistError(err error) bool {
	return errors.As(err, &userNotExistError{})
}
