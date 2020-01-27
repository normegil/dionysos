package stock

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/normegil/dionysos"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/model"
	"net/http"
)

type Controller struct {
	ItemDAO    ItemDAO
	ErrHandler error2.HTTPErrorHandler
}

func (c Controller) Routes() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/", c.loadAll)
	rt.Get("/{itemID}", c.load)
	return rt
}

func (c Controller) loadAll(w http.ResponseWriter, _ *http.Request) {
	items, err := c.ItemDAO.LoadAll()
	if nil != err {
		c.ErrHandler.Handle(w, fmt.Errorf("loading items: %w", err))
		return
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

type ItemDAO interface {
	LoadAll() ([]dionysos.Item, error)
}
