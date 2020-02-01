package database

import (
	"database/sql"
	"fmt"
	"github.com/gchaincl/dotsql"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos"
	"github.com/normegil/postgres"
	"github.com/rs/zerolog/log"
)

type ItemDAO struct {
	db      *sql.DB
	queries map[string]string
}

func (dao ItemDAO) NewItemDAO(db *sql.DB, owner string) (*ItemDAO, error) {
	sqlDir := pkger.Dir("sql/")
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
		var item dionysos.Item
		if err := rows.Scan(&item); nil != err {
			return nil, fmt.Errorf("scanning rows of item.'%s': %w", queryName, err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); nil != err {
		return nil, fmt.Errorf("parsing rows queried by item.'%s': %w", queryName, err)
	}
	return items, nil
}
