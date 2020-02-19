package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/postgres"
	"github.com/rs/zerolog/log"
)

type ItemDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB, owner string) (*ItemDAO, error) {
	queries := make(map[string]string)
	queries["Table-Existence"] = `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = 'item');`
	queries["Table-Create"] = `CREATE TABLE item ( id   uuid primary key, name varchar(300));`
	queries["Table-Set-Owner"] = `ALTER TABLE item OWNER TO $1;`
	err := postgres.CreateTable(db, postgres.TableInfos{
		Queries: queries,
		Owner:   owner,
	})
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &ItemDAO{db}, nil
}

func (dao *ItemDAO) LoadAll(opts model.CollectionOptions) ([]dionysos.Item, error) {
	rows, err := dao.db.Query(`SELECT * FROM item WHERE UPPER(item.name) LIKE UPPER($1) ORDER BY name LIMIT $2 OFFSET $3;`, toDatabaseFilter(opts.Filter), opts.Limit.Number(), opts.Offset.Number())
	if err != nil {
		return nil, fmt.Errorf("select all items %+v: %w", opts, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("select all: could not close rows")
		}
	}()

	items := make([]dionysos.Item, 0)
	for rows.Next() {
		var idStr string
		var name string
		if err := rows.Scan(&idStr, &name); nil != err {
			return nil, fmt.Errorf("scanning rows of items: %w", err)
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("cannot convert id '%s' to UUID : %w", idStr, err)
		}
		items = append(items, dionysos.Item{
			ID:   id,
			Name: name,
		})
	}

	if err := rows.Err(); nil != err {
		return nil, fmt.Errorf("parsing rows queried by items: %w", err)
	}
	return items, nil
}

func (dao ItemDAO) TotalNumberOfItem(filter string) (*model.Natural, error) {
	row := dao.db.QueryRow(`SELECT count(*) FROM item WHERE UPPER(item.name) LIKE UPPER($1);`, toDatabaseFilter(filter))
	var total int
	if err := row.Scan(&total); nil != err {
		return nil, fmt.Errorf("counting number of items: %w", err)
	}
	return model.NewNatural(total)
}

func (dao *ItemDAO) Insert(item dionysos.Item) error {
	if _, err := dao.db.Exec(`INSERT INTO item(id, name) VALUES (gen_random_uuid(), $1)`, item.Name); err != nil {
		return fmt.Errorf("inserting %+v: %w", item, err)
	}
	return nil
}

func (dao ItemDAO) Delete(itemID uuid.UUID) error {
	if _, err := dao.db.Exec(fmt.Sprintf("DELETE FROM item WHERE id = '%s'", itemID.String())); nil != err {
		return fmt.Errorf("deleting item '%s': %w", itemID.String(), err)
	}
	return nil
}
