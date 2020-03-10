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
	"strconv"
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
			items := generateItems(100)
			err := (&database.ItemDAO{Querier: tx}).InsertAll(items)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("WHEN query all items", func(t *testing.T) {
				resp := httptest.NewRecorder()
				lst.ServeHTTP(resp, test.Request(t, "GET", "/api/items", strings.NewReader("")))
				test.HandlerErrorResponse(t, resp)

				var body internalhttp.CollectionResponse
				test.FromJSONBody(t, resp, &body)

				t.Run("THEN filter is empty", func(t *testing.T) {
					if "" != body.Filter {
						t.Errorf("{expected:%s;got:%s}", "", body.Filter)
					}
				})

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

				t.Run("THEN number of items is the total number of items", func(t *testing.T) {
					if len(items) == body.NumberOfItems {
						t.Errorf("{expected:%d;got:%d}", len(items), body.NumberOfItems)
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

			t.Run("WHEN query filtered items", func(t *testing.T) {
				filter := "aaa"
				resp := httptest.NewRecorder()
				lst.ServeHTTP(resp, test.Request(t, "GET", "/api/items?filter="+filter, strings.NewReader("")))
				test.HandlerErrorResponse(t, resp)

				var body internalhttp.CollectionResponse
				test.FromJSONBody(t, resp, &body)

				t.Run("THEN filter is equal to requested filter", func(t *testing.T) {
					if "" != body.Filter {
						t.Errorf("{expected:%s;got:%s}", filter, body.Filter)
					}
				})

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

				filteredItems := make([]dionysos.Item, 0)
				for _, item := range items {
					if strings.Contains(item.Name, filter) {
						filteredItems = append(filteredItems, item)
					}
				}
				t.Run("THEN number of items is the total number of filtered items", func(t *testing.T) {
					if len(filteredItems) != body.NumberOfItems {
						t.Errorf("{expected:%d;got:%d}", len(items), body.NumberOfItems)
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

func generateItems(number int) []dionysos.Item {
	itemNamePrefixes := []string{
		"test-item-",
		"my-item-",
		"aaa-item-",
	}
	items := make([]dionysos.Item, 0)
	for i := 0; i < number; i++ {
		prefix := itemNamePrefixes[i%len(itemNamePrefixes)]
		items = append(items, dionysos.Item{
			ID:   uuid.New(),
			Name: prefix + strconv.Itoa(i),
		})
	}
	return items
}
