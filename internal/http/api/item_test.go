package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/normegil/dionysos"
	http2 "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestItemController(t *testing.T) {
	rt := chi.NewRouter()
	rt.Mount("/api/items", api.ItemController{ErrHandler: error2.HTTPErrorHandler{}}.Route())
	srv := httptest.NewServer(rt)
	defer srv.Close()

	t.Run("loadAll", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/items", srv.URL))
		if err != nil {
			t.Fatal(err)
		}
		var jsonItems json.RawMessage
		response := http2.CollectionResponse{
			Items: &jsonItems,
		}
		parseResponse(t, resp, &response)

		var items []dionysos.Item
		if err := json.Unmarshal(jsonItems, &items); nil != err {
			t.Fatal(fmt.Errorf("could not unmarshall '%+v': %w", items, err))
		}

		expectedItemNames := make([]string, 0)
		for i := 0; i < 10; i++ {
			expectedItemNames = append(expectedItemNames, "Item"+strconv.Itoa(i))
		}
		if len(items) != len(expectedItemNames) {
			t.Errorf("Wrong number of items {Expected:%d;Got:%d}", len(expectedItemNames), len(items))
		}
		for _, searched := range expectedItemNames {
			if !exist(searched, items) {
				t.Errorf("Expected item not found: %s", searched)
			}
		}
	})
}

func parseResponse(t testing.TB, resp *http.Response, result interface{}) {
	if http.StatusOK != resp.StatusCode && http.StatusNoContent != resp.StatusCode {
		t.Fatalf("Wrong status code {Got:%d;Expected:[%d,%d]}", resp.StatusCode, http.StatusOK, http.StatusNoContent)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(bodyBytes, result); nil != err {
		t.Fatal(fmt.Errorf("could not unmarshall '%s': %w", string(bodyBytes), err))
	}
}

func exist(searched string, items []dionysos.Item) bool {
	for _, item := range items {
		if searched == item.Name {
			return true
		}
	}
	return false
}
