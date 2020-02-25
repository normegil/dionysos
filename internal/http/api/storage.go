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

type StorageController struct {
	StorageDAO StorageDAO
	ErrHandler httperror.HTTPErrorHandler
}

func (c StorageController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/", c.loadAll)
	rt.Get("/{storageID}", c.load)
	rt.Put("/", c.save)
	rt.Delete("/{storageID}", c.delete)
	return rt
}

func (c StorageController) loadAll(w http.ResponseWriter, r *http.Request) {
	parameters := internalHTTP.QueryParameters{Params: r.URL.Query()}
	collectionOpts, err := internalHTTP.ToCollectionOptions(parameters)
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40010,
			Status: http.StatusBadRequest,
			Err:    err,
		})
		return
	}
	collectionOpts = c.toDefaultCollectionOptions(collectionOpts)

	storages, err := c.StorageDAO.LoadAll(*collectionOpts)
	if nil != err {
		c.ErrHandler.Handle(w, fmt.Errorf("loading storages: %w", err))
		return
	}

	nbStorages, err := c.StorageDAO.TotalNumberOfItem(collectionOpts.Filter)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("loading number of storages: %w", err))
		return
	}

	response := internalHTTP.CollectionResponse{
		Offset:        collectionOpts.Offset.Number(),
		Limit:         collectionOpts.Limit.Number(),
		NumberOfItems: nbStorages.Number(),
		Filter:        collectionOpts.Filter,
		Items:         storages,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		c.ErrHandler.Handle(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); nil != err {
		c.ErrHandler.Handle(w, err)
		return
	}
}

func (c StorageController) load(w http.ResponseWriter, _ *http.Request) {
	c.ErrHandler.Handle(w, httperror.HTTPError{
		Code:   50001,
		Status: http.StatusInternalServerError,
		Err:    errors.New("not implemented"),
	})
}

func (c StorageController) save(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("reading body: %w", err))
		return
	}

	var dto storageDTO
	if err := json.Unmarshal(bodyBytes, &dto); nil != err {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40011,
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("could not parse body '%s' into storage: %w", string(bodyBytes), err),
		})
		return
	}
	storage, err := dto.ToStorage()
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40011,
			Status: http.StatusBadRequest,
			Err:    err,
		})
	}

	created, err := c.StorageDAO.Save(storage)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("saving storage: %w", err))
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c StorageController) delete(w http.ResponseWriter, r *http.Request) {
	storageIDStr := chi.URLParam(r, "storageID")
	storageID, err := uuid.Parse(storageIDStr)
	if err != nil {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40010,
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("could not parse '%s' into uuid: %w", storageIDStr, err),
		})
		return
	}
	if err := c.StorageDAO.Delete(storageID); nil != err {
		c.ErrHandler.Handle(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c StorageController) toDefaultCollectionOptions(options *model.CollectionOptions) *model.CollectionOptions {
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

type StorageDAO interface {
	LoadAll(options model.CollectionOptions) ([]dionysos.Storage, error)
	Save(storage dionysos.Storage) (bool, error)
	TotalNumberOfItem(filter string) (*model.Natural, error)
	Delete(id uuid.UUID) error
}

type storageDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (dto storageDTO) ToStorage() (dionysos.Storage, error) {
	id := uuid.Nil
	if "" != dto.ID {
		var err error
		id, err = uuid.Parse(dto.ID)
		if err != nil {
			return dionysos.Storage{}, fmt.Errorf("parsing '%s': %w", dto.ID, err)
		}
	}
	return dionysos.Storage{
		ID:   id,
		Name: dto.Name,
	}, nil
}
