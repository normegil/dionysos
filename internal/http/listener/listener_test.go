package listener_test

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/dao/database"
	internalhttp "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	"github.com/normegil/dionysos/internal/http/listener"
	"github.com/normegil/dionysos/internal/tools/test"
	"github.com/normegil/postgres"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListener(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test")
	}

	lst, db := initTest(t)
	t.Run("GIVEN items exists", func(t *testing.T) {
		test.Transaction(t, db, func(tx *sql.Tx) {
			items := getTestItems()
			insertItems(t, tx, items)

			t.Run("WHEN query all items", func(t *testing.T) {
				resp := httptest.NewRecorder()
				lst.ServeHTTP(resp, test.Request(t, "GET", "/api/items", strings.NewReader("")))

				var body internalhttp.CollectionResponse
				test.FromJSONBody(t, resp, &body)

				t.Run("THEN collection offset is the starting offset", func(t *testing.T) {
					if body.Offset != 0 {
						t.Errorf("{expected:%d;got:%d}", 0, body.Offset)
					}
				})
				t.Run("THEN collection limit is the default", func(t *testing.T) {
					if body.Limit != api.DefaultLimit {
						t.Errorf("{expected:%d;got:%d}", api.DefaultLimit, body.Limit)
					}
				})

				respItems := body.Items.([]dionysos.Item)
				t.Run("THEN a limited number of items are returned", func(t *testing.T) {
					expected := len(items)
					if api.DefaultLimit < expected {
						expected = api.DefaultLimit
					}
					if len(respItems) != expected {
						t.Errorf("{expected:%d;got:%d}", len(items), len(respItems))
					}
				})
			})
		})
	})
}

func initTest(t testing.TB) (http.Handler, *sql.DB) {
	cfg, containerCfg := postgres.Test_Deploy(t)
	defer postgres.Test_RemoveContainer(t, containerCfg.Identifier)

	dbCfg := postgres.Configuration{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Database: "dionysos-listener-test-" + uuid.New().String(),
	}
	lst := listener.NewListener(listener.Configuration{
		APILogErrors: false,
		Database:     dbCfg,
	})
	defer func() {
		if err := lst.Close(); nil != err {
			t.Fatal(fmt.Errorf("closing listener: %w", err))
		}
	}()

	handler, err := lst.Load()
	if err != nil {
		t.Fatal(fmt.Errorf("load handler: %w", err))
	}

	db, err := postgres.New(dbCfg)
	if err != nil {
		t.Fatal(err)
	}

	return handler, db
}

func insertItems(t *testing.T, tx database.Querier, items []dionysos.Item) {
	dao := database.ItemDAO{Querier: tx}
	for _, item := range items {
		if err := dao.Insert(item); nil != err {
			t.Fatal(err)
		}
	}
}

func getTestItems() []dionysos.Item {
	items := make([]dionysos.Item, 0)
	items = append(items, dionysos.Item{
		ID:   uuid.New(),
		Name: "test-item-1",
	})
	items = append(items, dionysos.Item{
		ID:   uuid.New(),
		Name: "test-item-2",
	})
	items = append(items, dionysos.Item{
		ID:   uuid.New(),
		Name: "test-item-3",
	})
	return items
}
