package versions

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/database"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/postgres"
)

type SchemaCreation struct {
	DB *sql.DB
}

func (v SchemaCreation) Apply() error {
	if err := v.createTables(); err != nil {
		return fmt.Errorf("creating tables: %w", err)
	}
	if err := v.insertData(); err != nil {
		return fmt.Errorf("inserting datas: %w", err)
	}
	return nil
}

func (v SchemaCreation) insertData() error {
	defaultRoles := []security.Role{
		{Name: "root"},
		{Name: "user"},
	}
	for _, role := range defaultRoles {
		if _, err := v.DB.Exec(`INSERT INTO role (id, name) VALUES (gen_random_uuid(), $1)`, role.Name); nil != err {
			return fmt.Errorf("inserting '%s' role: %w", err)
		}
	}

	dao := database.CasbinDAO{DB: v.DB}
	if err := dao.InsertMultiple(defaultPolicies()); nil != err {
		return fmt.Errorf("inserting default rules: %w", err)
	}

	return nil
}

func defaultPolicies() []model.CasbinRule {
	rules := make([]model.CasbinRule, 0)
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, item, read"})
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, item, write"})
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, storage, read"})
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, storage, write"})
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "none, item, read"})
	rules = append(rules, model.CasbinRule{ID: uuid.Nil, Type: "p", Value: "none, storage, read"})
	return rules
}

func (v SchemaCreation) createTables() error {
	row := v.DB.QueryRow(`SELECT pg_catalog.pg_get_userbyid(d.datdba) as "Owner" FROM pg_catalog.pg_database d WHERE d.datname = current_database();`)
	var owner string
	if err := row.Scan(&owner); nil != err {
		return fmt.Errorf("load database owner: %w", err)
	}

	tableExistence := `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = '%s');`
	tableSetOwner := `ALTER TABLE "%s" OWNER TO $1;`

	itemTableName := "item"
	err := postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, itemTableName),
			"Table-Create":    `CREATE TABLE item ( id uuid primary key, name varchar(300));`,
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
			"Table-Create":    `CREATE TABLE storage ( id uuid primary key, name varchar(300));`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, storageTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", storageTableName, err)
	}

	roleTableName := "role"
	err = postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, roleTableName),
			"Table-Create":    `CREATE TABLE role ( id uuid primary key, name varchar(300));`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, roleTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", roleTableName, err)
	}

	userTableName := "user"
	err = postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, userTableName),
			"Table-Create": `CREATE TABLE "user" (
				id uuid primary key,
				name varchar(300) unique,
				hash bytea,
				algorithmID uuid,
				roleID uuid REFERENCES role(id) ON DELETE RESTRICT);`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, userTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", userTableName, err)
	}

	policyTableName := "policy"
	err = postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, policyTableName),
			"Table-Create": `CREATE TABLE "policy" (
				id varchar(300) primary key,
				ptype varchar(300),
				value varchar(300)
			);`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, policyTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", policyTableName, err)
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
