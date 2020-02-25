package versions

import (
	"database/sql"
	"fmt"
	"github.com/normegil/postgres"
)

type SchemaCreation struct {
	DB *sql.DB
}

func (v SchemaCreation) Apply() error {
	row := v.DB.QueryRow(`SELECT pg_catalog.pg_get_userbyid(d.datdba) as "Owner" FROM pg_catalog.pg_database d WHERE d.datname = current_database();`)
	var owner string
	if err := row.Scan(&owner); nil != err {
		return fmt.Errorf("load database owner: %w", err)
	}

	tableExistence := `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = '%s');`
	tableSetOwner := `ALTER TABLE %s OWNER TO $1;`

	itemTableName := "item"
	err := postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, itemTableName),
			"Table-Create":    `CREATE TABLE item ( id   uuid primary key, name varchar(300));`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, itemTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", itemTableName, err)
	}

	storageTableName := "storage"
	err = postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, storageTableName),
			"Table-Create":    `CREATE TABLE storage ( id   uuid primary key, name varchar(300));`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, storageTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", storageTableName, err)
	}

	userTableName := "user"
	err = postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, userTableName),
			"Table-Create": `CREATE TABLE "user" (
				id uuid primary key,
				name varchar(300) unique,
				hash bytea,
				algorithmID uuid);`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, userTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", userTableName, err)
	}

	return nil
}

func (v SchemaCreation) Rollback() error {
	dropTableQuery := "DROP TABLE %s"

	userTableName := "user"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, userTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", userTableName, err)
	}

	storageTableName := "storage"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, storageTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", storageTableName, err)
	}

	itemTableName := "item"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, itemTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", itemTableName, err)
	}

	return nil
}
