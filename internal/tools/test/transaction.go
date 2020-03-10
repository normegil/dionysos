package test

import (
	"database/sql"
	"testing"
)

func Transaction(t testing.TB, db *sql.DB, toCall func(tx *sql.Tx)) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	toCall(tx)
	if err = tx.Rollback(); nil != err {
		t.Fatal(err)
	}
}
