package database

import (
	"database/sql"
	"fmt"
	"github.com/gchaincl/dotsql"
	"github.com/google/uuid"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos"
	"github.com/normegil/postgres"
	"github.com/rs/zerolog/log"
)

type ItemDAO struct {
	db      *sql.DB
	queries map[string]string
}

func NewItemDAO(db *sql.DB, owner string) (*ItemDAO, error) {
	sqlDir := pkger.Dir("/internal/database/sql")
	fileName := "item.sql"
	file, err := sqlDir.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("loading file %s: %w", fileName, err)
	}
	defer func() {
		if err := file.Close(); nil != err {
			log.Error().Err(err).Msgf("could not close file %s", fileName)
		}
	}()

	itemSql, err := dotsql.Load(file)
	if err != nil {
		return nil, fmt.Errorf("parsing SQL file %s: %w", fileName, err)
	}

	queries := itemSql.QueryMap()
	err = postgres.CreateTable(db, postgres.TableInfos{
		Queries: queries,
		Owner:   owner,
	})
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}
	return &ItemDAO{db, queries}, nil
}

func (dao *ItemDAO) LoadAll() ([]dionysos.Item, error) {
	queryName := "Select-All"
	rows, err := dao.db.Query(dao.queries[queryName])
	if err != nil {
		return nil, fmt.Errorf("query failed item.'%s': %w", queryName, err)
	}
	defer func() {
		if err := rows.Close(); nil != err {
			log.Error().Err(err).Msgf("could not close rows of item.'%s'", queryName)
		}
	}()

	items := make([]dionysos.Item, 0)
	for rows.Next() {
		var idStr string
		var name string
		if err := rows.Scan(&idStr, &name); nil != err {
			return nil, fmt.Errorf("scanning rows of item.'%s': %w", queryName, err)
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
		return nil, fmt.Errorf("parsing rows queried by item.'%s': %w", queryName, err)
	}
	return items, nil
}

func (dao *ItemDAO) Insert(item dionysos.Item) error {
	queryName := "Insert"
	if _, err := dao.db.Exec(dao.queries[queryName], item.Name); err != nil {
		return fmt.Errorf("inserting %+v: %w", item, err)
	}
	return nil
}
