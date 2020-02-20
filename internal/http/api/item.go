package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/model"
	"io/ioutil"
	"net/http"
)

type ItemController struct {
	ItemDAO    ItemDAO
	ErrHandler httperror.HTTPErrorHandler
}

func (c ItemController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/", c.loadAll)
	rt.Get("/{itemID}", c.load)
	rt.Put("/", c.save)
	rt.Delete("/{itemID}", c.delete)
	return rt
}

func (c ItemController) loadAll(w http.ResponseWriter, r *http.Request) {
	parameters := internalHTTP.QueryParameters{Params: r.URL.Query()}
	collectionOpts, err := internalHTTP.ToCollectionOptions(parameters)
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   400,
			Status: 40010,
			Err:    err,
		})
		return
	}
	collectionOpts = c.toDefaultCollectionOptions(collectionOpts)

	items, err := c.ItemDAO.LoadAll(*collectionOpts)
	if nil != err {
		c.ErrHandler.Handle(w, fmt.Errorf("loading items: %w", err))
		return
	}

	nbItems, err := c.ItemDAO.TotalNumberOfItem(collectionOpts.Filter)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("loading number of items: %w", err))
		return
	}

	response := internalHTTP.CollectionResponse{
		Offset:        collectionOpts.Offset.Number(),
		Limit:         collectionOpts.Limit.Number(),
		NumberOfItems: nbItems.Number(),
		Filter:        collectionOpts.Filter,
		Items:         items,
	}

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

func (c ItemController) load(w http.ResponseWriter, _ *http.Request) {
	c.ErrHandler.Handle(w, httperror.HTTPError{
		Code:   50001,
		Status: http.StatusInternalServerError,
		Err:    errors.New("not implemented"),
	})
}

func (c ItemController) save(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("reading body: %w", err))
		return
	}

	var dto itemDTO
	if err := json.Unmarshal(bodyBytes, &dto); nil != err {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40011,
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("could not parse body '%s' into item: %w", string(bodyBytes), err),
		})
		return
	}

	item, err := dto.ToItem()
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40011,
			Status: http.StatusBadRequest,
			Err:    err,
		})
	}

	created, err := c.ItemDAO.Save(item)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("saving item: %w", err))
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c ItemController) delete(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40010,
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("could not parse '%s' into uuid: %w", itemIDStr, err),
		})
		return
	}
	if err := c.ItemDAO.Delete(itemID); nil != err {
		c.ErrHandler.Handle(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c ItemController) toDefaultCollectionOptions(options *model.CollectionOptions) *model.CollectionOptions {
	if 0 != options.Limit.Number() {
		return options
	}
	newLimit, err := model.NewNatural(50)
	if err != nil {
		panic(fmt.Errorf("should not fail: %w", err))
	}
	return &model.CollectionOptions{
		Limit:  *newLimit,
		Offset: options.Offset,
		Filter: options.Filter,
	}
}

type ItemDAO interface {
	LoadAll(options model.CollectionOptions) ([]dionysos.Item, error)
	Save(item dionysos.Item) (bool, error)
	TotalNumberOfItem(filter string) (*model.Natural, error)
	Delete(id uuid.UUID) error
}

type itemDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (dto itemDTO) ToItem() (dionysos.Item, error) {
	id := uuid.Nil
	if "" != dto.ID {
		var err error
		id, err = uuid.Parse(dto.ID)
		if err != nil {
			return dionysos.Item{}, fmt.Errorf("parsing '%s': %w", dto.ID, err)
		}
	}
	return dionysos.Item{
		ID:   id,
		Name: dto.Name,
	}, nil
}
