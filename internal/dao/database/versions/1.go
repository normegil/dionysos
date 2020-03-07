package versions

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/configuration"
	"github.com/normegil/dionysos/internal/dao/database"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
	roleDAO := database.RoleDAO{DB: v.DB}

	rootRole := security.Role{
		Name: "root",
	}

	defaultRoles := []security.Role{
		rootRole,
		{Name: "user"},
	}
	for _, role := range defaultRoles {
		if err := roleDAO.Insert(role); nil != err {
			return fmt.Errorf("inserting '%s' role: %w", role.Name, err)
		}
	}

	dao := database.CasbinDAO{DB: v.DB}
	if err := dao.InsertMultiple(defaultPolicies()); nil != err {
		return fmt.Errorf("inserting default rules: %w", err)
	}

	loadedRootRole, err := roleDAO.LoadByName(rootRole.Name)
	if err != nil {
		return fmt.Errorf("could not find role '%s': %w", rootRole.Name, err)
	}

	userDAO := database.UserDAO{DB: v.DB}
	user, err := security.NewUser("root", "root", *loadedRootRole)
	if err != nil {
		return fmt.Errorf("creating 'admin' user: %w", err)
	}
	if err := userDAO.Insert(*user); nil != err {
		return fmt.Errorf("inserting '%s' user: %w", user.Name, err)
	}

	if err := v.insertDummyData(v.DB); nil != err {
		return fmt.Errorf("inserting dummy data: %w", err)
	}

	return nil
}

func defaultPolicies() []security.CasbinRule {
	rules := make([]security.CasbinRule, 0)
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, item, read"})
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, item, write"})
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, storage, read"})
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "user, storage, write"})
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "none, item, read"})
	rules = append(rules, security.CasbinRule{ID: uuid.Nil, Type: "p", Value: "none, storage, read"})
	return rules
}

//nolint:funlen // Function is highly repetitive with low complexity
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
			"Table-Create":    `CREATE TABLE role ( id uuid primary key, name varchar(300) unique);`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, roleTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", roleTableName, err)
	}

	const userTableName = "user"
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
				type varchar(300),
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

	policyTableName := "policy"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, policyTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", policyTableName, err)
	}

	userTableName := "user"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, userTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", userTableName, err)
	}

	roleTableName := "role"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, roleTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", roleTableName, err)
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

func (v SchemaCreation) insertDummyData(db *sql.DB) error {
	if viper.GetBool(configuration.KeyDummyData.Name) {
		storageDAO := &database.StorageDAO{DB: db}
		itemDAO := &database.ItemDAO{DB: db}
		userDAO := &database.UserDAO{DB: db}

		log.Info().Msg("insert dummy data")
		const dummyItemNb = 50

		rolename := "user"
		userrole, err := database.RoleDAO{DB: v.DB}.LoadByName(rolename)
		if err != nil {
			return fmt.Errorf("load '%s' role: %w", rolename, err)
		}

		user, err := security.NewUser("user", "user", *userrole)
		if err != nil {
			return fmt.Errorf("creating 'user' user: %w", err)
		}
		if err := userDAO.Insert(*user); nil != err {
			return fmt.Errorf("inserting %s: %w", user.Name, err)
		}

		for i := 0; i < dummyItemNb; i++ {
			storage := dionysos.Storage{Name: gofakeit.Company()}
			if err := storageDAO.Insert(storage); nil != err {
				return fmt.Errorf("inserting %+v: %w", storage, err)
			}
		}

		for i := 0; i < dummyItemNb; i++ {
			item := dionysos.Item{Name: gofakeit.BeerName()}
			if err := itemDAO.Insert(item); nil != err {
				return fmt.Errorf("inserting %+v: %w", item, err)
			}
		}
	}
	return nil
}
