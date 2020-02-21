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

type StorageDAO struct {
	db *sql.DB
}

func NewStorageDAO(db *sql.DB, owner string) (*StorageDAO, error) {
	queries := make(map[string]string)
	queries["Table-Existence"] = `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = 'storage');`
	queries["Table-Create"] = `CREATE TABLE storage ( id   uuid primary key, name varchar(300));`
	queries["Table-Set-Owner"] = `ALTER TABLE storage OWNER TO $1;`
	err := postgres.CreateTable(db, postgres.TableInfos{
		Queries: queries,
		Owner:   owner,
	})
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &StorageDAO{db}, nil
}

func (dao *StorageDAO) LoadAll(opts model.CollectionOptions) ([]dionysos.Storage, error) {
	rows, err := dao.db.Query(`SELECT * FROM storage WHERE UPPER(storage.name) LIKE UPPER($1) ORDER BY name LIMIT $2 OFFSET $3;`, toDatabaseFilter(opts.Filter), opts.Limit.Number(), opts.Offset.Number())
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
	row := dao.db.QueryRow(`SELECT count(*) FROM storage WHERE UPPER(storage.name) LIKE UPPER($1);`, toDatabaseFilter(filter))
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
	if _, err := dao.db.Exec(`INSERT INTO storage(id, name) VALUES (gen_random_uuid(), $1)`, storage.Name); err != nil {
		return fmt.Errorf("inserting %+v: %w", storage, err)
	}
	return nil
}

func (dao *StorageDAO) Update(storage dionysos.Storage) error {
	if _, err := dao.db.Exec(`UPDATE storage SET name = $2 WHERE id = $1`, storage.ID, storage.Name); err != nil {
		return fmt.Errorf("updating %+v: %w", storage, err)
	}
	return nil
}

func (dao StorageDAO) Delete(storageID uuid.UUID) error {
	if _, err := dao.db.Exec(fmt.Sprintf("DELETE FROM storage WHERE id = '%s'", storageID.String())); nil != err {
		return fmt.Errorf("deleting storage '%s': %w", storageID.String(), err)
	}
	return nil
}
