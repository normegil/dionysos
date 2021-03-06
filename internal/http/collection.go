package http

import (
	"fmt"
	"github.com/normegil/dionysos/internal/model"
	"strconv"
)

type CollectionResponse struct {
	Offset        int         `json:"offset"`
	Limit         int         `json:"limit"`
	Filter        string      `json:"filter"`
	NumberOfItems int         `json:"totalSize"`
	Items         interface{} `json:"items"`
}

func ToCollectionOptions(parameters QueryParameters) (*model.CollectionOptions, error) {
	limit, err := ToNatural(parameters, "limit")
	if err != nil {
		return nil, err
	}

	offset, err := ToNatural(parameters, "offset")
	if err != nil {
		return nil, err
	}

	filterKey := "filter"
	filter, err := parameters.ExtractSingle(filterKey)
	if err != nil {
		return nil, fmt.Errorf("extracting '%s': %w", filterKey, err)
	}

	return &model.CollectionOptions{
		Limit:  *limit,
		Offset: *offset,
		Filter: filter,
	}, nil
}

func ToNatural(parameters QueryParameters, key string) (*model.Natural, error) {
	limitStr, err := parameters.ExtractSingle(key)
	if err != nil {
		return nil, err
	}
	if "" == limitStr {
		return model.NewNatural(0)
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}
	return model.NewNatural(limit)
}
