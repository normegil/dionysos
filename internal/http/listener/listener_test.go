package listener_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/dao/database"
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

//nolint:funlen,gocognit
func TestListener(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test")
	}

	lst, db, containerCfg := initTest(t)
	defer postgres.Test_RemoveContainer(t, containerCfg.Identifier)
	t.Run("GIVEN items exists", func(t *testing.T) {
		items := generateItems(100)
		dao := &database.ItemDAO{Querier: db}
		err := dao.InsertAll(items)
		if err != nil {
			t.Fatal(err)
		}
		defer test.Clear(t, []test.Clearer{dao})

		t.Run("WHEN query all items", func(t *testing.T) {
			resp := httptest.NewRecorder()
			lst.ServeHTTP(resp, test.Request(t, "GET", "/api/items", strings.NewReader("")))
			respBody := test.ReadResponse(t, resp)
			test.HandlerErrorResponse(t, resp.Code, respBody)

			var body CollectionResponse
			test.FromJSONBody(t, respBody, &body)

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
				if len(items) != body.NumberOfItems {
					t.Errorf("{expected:%d;got:%d}", len(items), body.NumberOfItems)
				}
			})

			var respItems []dionysos.Item
			if err := json.Unmarshal(body.Items, &respItems); nil != err {
				t.Fatal(err)
			}
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
			respBody := test.ReadResponse(t, resp)
			test.HandlerErrorResponse(t, resp.Code, respBody)

			var body CollectionResponse
			test.FromJSONBody(t, respBody, &body)

			t.Run("THEN filter is equal to requested filter", func(t *testing.T) {
				if filter != body.Filter {
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

			var respItems []dionysos.Item
			if err := json.Unmarshal(body.Items, &respItems); nil != err {
				t.Fatal(err)
			}
			t.Run("THEN a limited number of items are returned", func(t *testing.T) {
				expected := len(filteredItems)
				if api.DefaultLimit < expected {
					expected = api.DefaultLimit
				}
				if len(respItems) != expected {
					t.Errorf("{expected:%d;got:%d}", expected, len(respItems))
				}
			})
		})
	})
}

func initTest(t testing.TB) (http.Handler, *sql.DB, postgres.ContainerConfiguration) {
	cfg, containerCfg := postgres.Test_Deploy(t)
	dbCfg := postgres.Configuration{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Database: "dionysos_listener_test_" + strings.ReplaceAll(uuid.New().String(), "-", "_"),
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

	return handler, db, containerCfg
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

type CollectionResponse struct {
	Offset        int             `json:"offset"`
	Limit         int             `json:"limit"`
	Filter        string          `json:"filter"`
	NumberOfItems int             `json:"totalSize"`
	Items         json.RawMessage `json:"items"`
}
