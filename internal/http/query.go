package http

import (
	"fmt"
	"net/url"
)

type QueryParameters struct {
	Params url.Values
}

func (q QueryParameters) ExtractSingle(key string) (string, error) {
	parameter := q.Params[key]
	if nil == parameter {
		return "", nil
	}
	if len(parameter) > 1 {
		return "", fmt.Errorf("more than one parameter defined for %s (found %d)", key, len(parameter))
	}
	return parameter[0], nil
}
