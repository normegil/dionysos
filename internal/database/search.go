package database

import (
	"github.com/normegil/dionysos/internal/model"
)

type Searcheable interface {
	Search(parameters model.SearchParameters) (*model.SearchResult, error)
}

type SearchDAO struct {
	Searchables []Searcheable
}

func (d SearchDAO) Search(searchParameters model.SearchParameters) ([]model.SearchResult, error) {
	results := make([]model.SearchResult, 0)
	for _, searchable := range d.Searchables {
		result, err := searchable.Search(searchParameters)
		if err != nil {
			return nil, err
		}
		if 0 != len(result.Results) {
			results = append(results, *result)
		}
	}
	return results, nil
}
