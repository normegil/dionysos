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

type StorageDAO struct {
	DB *sql.DB
}

func (dao StorageDAO) Resource() model.Resource {
	return model.ResourceStorage
}

func (dao *StorageDAO) LoadAll(opts model.CollectionOptions) ([]dionysos.Storage, error) {
	rows, err := dao.DB.Query(`SELECT * FROM storage WHERE UPPER(storage.name) LIKE UPPER($1) ORDER BY name LIMIT $2 OFFSET $3;`, toDatabaseFilter(opts.Filter), opts.Limit.Number(), opts.Offset.Number())
	if err != nil {
		return nil, fmt.Errorf("select all storages %+v: %w", opts, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("select all: could not close rows")
		}
	}()

	storages := make([]dionysos.Storage, 0)
	for rows.Next() {
		var idStr string
		var name string
		if err := rows.Scan(&idStr, &name); nil != err {
			return nil, fmt.Errorf("scanning rows of storages: %w", err)
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("cannot convert id '%s' to UUID : %w", idStr, err)
		}
		storages = append(storages, dionysos.Storage{
			ID:   id,
			Name: name,
		})
	}

	if err := rows.Err(); nil != err {
		return nil, fmt.Errorf("parsing rows queried by storages: %w", err)
	}
	return storages, nil
}

func (dao StorageDAO) TotalNumberOfItem(filter string) (*model.Natural, error) {
	row := dao.DB.QueryRow(`SELECT count(*) FROM storage WHERE UPPER(storage.name) LIKE UPPER($1);`, toDatabaseFilter(filter))
	var total int
	if err := row.Scan(&total); nil != err {
		return nil, fmt.Errorf("counting number of storages: %w", err)
	}
	return model.NewNatural(total)
}

func (dao StorageDAO) Save(storage dionysos.Storage) (bool, error) {
	if storage.ID == uuid.Nil {
		err := dao.Insert(storage)
		return true, err
	}
	err := dao.Update(storage)
	return false, err
}

func (dao *StorageDAO) Insert(storage dionysos.Storage) error {
	if _, err := dao.DB.Exec(`INSERT INTO storage(id, name) VALUES (gen_random_uuid(), $1)`, storage.Name); err != nil {
		return fmt.Errorf("inserting %+v: %w", storage, err)
	}
	return nil
}

func (dao *StorageDAO) Update(storage dionysos.Storage) error {
	if _, err := dao.DB.Exec(`UPDATE storage SET name = $2 WHERE id = $1`, storage.ID, storage.Name); err != nil {
		return fmt.Errorf("updating %+v: %w", storage, err)
	}
	return nil
}

func (dao StorageDAO) Delete(storageID uuid.UUID) error {
	if _, err := dao.DB.Exec(fmt.Sprintf("DELETE FROM storage WHERE id = '%s'", storageID.String())); nil != err {
		return fmt.Errorf("deleting storage '%s': %w", storageID.String(), err)
	}
	return nil
}

func (dao StorageDAO) Search(params model.SearchParameters) (*model.SearchResult, error) {
	query := `SELECT searchIndex.id, searchIndex.name
				FROM (
					 SELECT id, name, setweight(to_tsvector(name), 'A') as document
					 FROM storage
				 ) searchIndex`
	if 0 != len(params.Searches) {
		query += " WHERE "
	}
	for i, _ := range params.Searches {
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
	rows, err := dao.DB.Query(query, searches...)
	if err != nil {
		return nil, fmt.Errorf("search for storages with parameter %+v: %w", params, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("search storages: could not close rows")
		}
	}()

	storages, err := dao.fromRows(rows)
	storagesInt := make([]interface{}, 0)
	for _, storage := range storages {
		storagesInt = append(storagesInt, storage)
	}
	if nil != err {
		return nil, fmt.Errorf("map search results to storages: %w", err)
	}
	return &model.SearchResult{
		Type:    "Storage",
		Results: storagesInt,
	}, nil
}

func (dao StorageDAO) fromRows(rows *sql.Rows) ([]dionysos.Storage, error) {
	items := make([]dionysos.Storage, 0)
	for rows.Next() {
		var idStr string
		var name string
		if err := rows.Scan(&idStr, &name); nil != err {
			return nil, fmt.Errorf("scanning rows of storage: %w", err)
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("cannot convert storage id '%s' to UUID : %w", idStr, err)
		}
		items = append(items, dionysos.Storage{
			ID:   id,
			Name: name,
		})
	}

	if err := rows.Err(); nil != err {
		return nil, fmt.Errorf("parsing rows queried by storages: %w", err)
	}

	return items, nil
}
