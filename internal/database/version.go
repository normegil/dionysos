package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type VersionDAO struct {
	DB *sql.DB
}

func (d VersionDAO) CurrentVersion() (int, error) {
	row := d.DB.QueryRow("SELECT version FROM version ORDER BY modificationTime DESC LIMIT 1")
	var version int
	if err := row.Scan(&version); err != nil {
		if d.errIsTableNotExist(err) {
			return -1, nil
		}
		return -1, fmt.Errorf("could not get current version: %w", err)
	}
	return version, nil
}

func (d VersionDAO) errIsTableNotExist(err error) bool {
	return strings.Contains(err.Error(), "not exist")
}

func (d VersionDAO) Insert(version int) error {
	if _, err := d.DB.Exec(`INSERT INTO 
    		version(id, version, modificationTime)
    		VALUES (gen_random_uuid(), $1, now())`, version); nil != err {
		return fmt.Errorf("could not insert version %d: %w", version, err)
	}
	return nil
}
