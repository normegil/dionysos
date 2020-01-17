package stock

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/model"
	"net/http"
	"strconv"
)

type Controller struct {
	ErrHandler error2.HTTPErrorHandler
}

func (c Controller) Routes() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/", c.loadAll)
	rt.Get("/{itemID}", c.load)
	return rt
}

func (c Controller) loadAll(w http.ResponseWriter, _ *http.Request) {
	loadedItems := make([]dionysos.Item, 0)
	for i := 0; i < 10; i++ {
		loadedItems = append(loadedItems, dionysos.Item{
			ID:   uuid.New(),
			Name: "Item" + strconv.Itoa(i),
		})
	}

	items := make([]interface{}, 0)
	for _, item := range loadedItems {
		items = append(items, item)
	}

	response := model.CollectionResponse{Items: items}

	bytes, err := json.Marshal(response)
	if err != nil {
		c.ErrHandler.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); nil != err {
		c.ErrHandler.Handle(w, err)
		return
	}
}

func (c Controller) load(w http.ResponseWriter, _ *http.Request) {
	c.ErrHandler.Handle(w, error2.HTTPError{
		Code:   50001,
		Status: http.StatusInternalServerError,
		Err:    errors.New("not implemented"),
	})
}
