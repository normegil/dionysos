package database

import (
	"fmt"
	"github.com/normegil/dionysos/internal/model"
	"github.com/normegil/dionysos/internal/security"
)

type Searcheable interface {
	Resource() model.Resource
	Search(parameters model.SearchParameters) (*model.SearchResult, error)
}

type SearchDAO struct {
	Searchables []Searcheable
	Authorizer  security.Authorizer
}

func (d SearchDAO) Search(role security.Role, searchParameters model.SearchParameters) ([]model.SearchResult, error) {
	results := make([]model.SearchResult, 0)
	for _, searchable := range d.Searchables {
		authorized, err := d.Authorizer.IsAuthorized(role, searchable.Resource(), model.ActionRead)
		if err != nil {
			return nil, fmt.Errorf("checking authorization for '%s': %w", role.Name, err)
		}
		if authorized {
			result, err := searchable.Search(searchParameters)
			if err != nil {
				return nil, err
			}
			if 0 != len(result.Results) {
				results = append(results, *result)
			}
		}
	}
	return results, nil
}
