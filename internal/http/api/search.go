package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/model"
	"io/ioutil"
	"net/http"
	"strings"
)

type SearchController struct {
	DAO        SearchDAO
	ErrHandler httperror.HTTPErrorHandler
}

func (c SearchController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Put("/", c.search)
	return rt
}

func (c SearchController) search(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("reading body: %w", err))
		return
	}

	var dto searchDTO
	if err := json.Unmarshal(bodyBytes, &dto); nil != err {
		c.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   40011,
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("could not parse body '%s' into item: %w", string(bodyBytes), err),
		})
		return
	}

	params := dto.ToSearchParameters()
	results, err := c.DAO.Search(params)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("searching: %w", err))
		return
	}

	marshalled, err := json.Marshal(searchResponse{Results: results})
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("marshalling response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(marshalled); err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("write response body: %w", err))
		return
	}
}

type SearchDAO interface {
	Search(searchParameters model.SearchParameters) ([]model.SearchResult, error)
}

type searchDTO struct {
	Search string `json:"search"`
}

type searchResponse struct {
	Results []model.SearchResult `json:"results"`
}

func (d searchDTO) ToSearchParameters() model.SearchParameters {
	trimmedSearch := strings.TrimSpace(d.Search)
	splittedSearch := strings.Split(trimmedSearch, " ")
	searches := make([]string, 0)
	for _, search := range splittedSearch {
		searches = append(searches, strings.TrimSpace(search))
	}
	return model.SearchParameters{Searches: searches}
}
