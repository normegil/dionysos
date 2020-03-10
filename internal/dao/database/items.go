//nolint:dupl // Dummy model trigger duplicates warning
package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/model"
	"github.com/rs/zerolog/log"
	"strconv"
)

type ItemDAO struct {
	Querier Querier
}

func (dao ItemDAO) Resource() model.Resource {
	return model.ResourceItem
}

func (dao *ItemDAO) LoadAll(opts model.CollectionOptions) ([]dionysos.Item, error) {
	rows, err := dao.Querier.Query(`SELECT * FROM item WHERE UPPER(item.name) LIKE UPPER($1) ORDER BY name LIMIT $2 OFFSET $3;`, toDatabaseFilter(opts.Filter), opts.Limit.Number(), opts.Offset.Number())
	if err != nil {
		return nil, fmt.Errorf("select all items %+v: %w", opts, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("select all: could not close rows")
		}
	}()

	return dao.fromRows(rows)
}

func (dao ItemDAO) TotalNumberOfItem(filter string) (*model.Natural, error) {
	row := dao.Querier.QueryRow(`SELECT count(*) FROM item WHERE UPPER(item.name) LIKE UPPER($1);`, toDatabaseFilter(filter))
	var total int
	if err := row.Scan(&total); nil != err {
		return nil, fmt.Errorf("counting number of items: %w", err)
	}
	return model.NewNatural(total)
}

func (dao ItemDAO) Save(item dionysos.Item) (bool, error) {
	if item.ID == uuid.Nil {
		err := dao.Insert(item)
		return true, err
	}
	err := dao.Update(item)
	return false, err
}

func (dao *ItemDAO) Insert(item dionysos.Item) error {
	if _, err := dao.Querier.Exec(`INSERT INTO item(id, name) VALUES (gen_random_uuid(), $1)`, item.Name); err != nil {
		return fmt.Errorf("inserting %+v: %w", item, err)
	}
	return nil
}

func (dao *ItemDAO) Update(item dionysos.Item) error {
	if _, err := dao.Querier.Exec(`UPDATE item SET name = $2 WHERE id = $1`, item.ID, item.Name); err != nil {
		return fmt.Errorf("updating %+v: %w", item, err)
	}
	return nil
}

func (dao ItemDAO) Delete(itemID fmt.Stringer) error {
	if _, err := dao.Querier.Exec("DELETE FROM item WHERE id = ?", itemID.String()); nil != err {
		return fmt.Errorf("deleting item '%s': %w", itemID.String(), err)
	}
	return nil
}

func (dao ItemDAO) Search(params model.SearchParameters) (*model.SearchResult, error) {
	query := `SELECT searchIndex.id, searchIndex.name
				FROM (
					 SELECT id, name, setweight(to_tsvector(name), 'A') as document
					 FROM item
				 ) searchIndex`
	if 0 != len(params.Searches) {
		query += " WHERE "
	}
	for i := range params.Searches {
		if i != 0 {
			query += " OR "
		}
		query += "searchIndex.document @@ to_tsquery($" + strconv.Itoa(2*i+1) + ") OR searchIndex.document @@ $" + strconv.Itoa(2*i+2)
	}

	searches := make([]interface{}, 0)
	for _, search := range params.Searches {
		if "" == search {
			searches = append(append(searches, "*"), "*")
		} else {
			modifiedSearchTerm := search + ":*"
			searches = append(append(searches, modifiedSearchTerm), modifiedSearchTerm)
		}
	}
	rows, err := dao.Querier.Query(query, searches...)
	if err != nil {
		return nil, fmt.Errorf("search for items with parameters %+v: %w", params, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("search items: could not close rows")
		}
	}()

	items, err := dao.fromRows(rows)
	interfaces := make([]interface{}, 0)
	for _, item := range items {
		interfaces = append(interfaces, item)
	}
	if nil != err {
		return nil, fmt.Errorf("map search results to items: %w", err)
	}
	return &model.SearchResult{
		Type:    "Item",
		Results: interfaces,
	}, nil
}

func (dao ItemDAO) fromRows(rows *sql.Rows) ([]dionysos.Item, error) {
	items := make([]dionysos.Item, 0)
	for rows.Next() {
		var idStr string
		var name string
		if err := rows.Scan(&idStr, &name); nil != err {
			return nil, fmt.Errorf("scanning rows of items: %w", err)
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("cannot convert item id '%s' to UUID : %w", idStr, err)
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
