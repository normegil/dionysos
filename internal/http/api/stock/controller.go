package stock

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	error2 "github.com/normegil/dionysos/internal/http/error"
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
	items := make([]dionysos.Item, 0)
	for i := 0; i < 10; i++ {
		items = append(items, dionysos.Item{
			ID:   uuid.New(),
			Name: "Item" + strconv.Itoa(i),
		})
	}

	bytes, err := json.Marshal(items)
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
