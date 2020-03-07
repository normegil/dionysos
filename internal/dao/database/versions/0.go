package versions

import (
	"database/sql"
	"fmt"
	"github.com/normegil/postgres"
)

type VersionCreation struct {
	DB *sql.DB
}

func (v VersionCreation) Apply() error {
	row := v.DB.QueryRow(`SELECT pg_catalog.pg_get_userbyid(d.datdba) as "Owner" FROM pg_catalog.pg_database d WHERE d.datname = current_database();`)
	var owner string
	if err := row.Scan(&owner); nil != err {
		return fmt.Errorf("load database owner: %w", err)
	}

	tableExistence := `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = '%s');`
	tableSetOwner := `ALTER TABLE %s OWNER TO $1;`

	versionTableName := "version"
	err := postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, versionTableName),
			"Table-Create":    `CREATE TABLE version (id uuid primary key, version integer, modificationTime timestamp)`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, versionTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", versionTableName, err)
	}

	return nil
}

func (v VersionCreation) Rollback() error {
	return fmt.Errorf("cannot rollback this change")
}
