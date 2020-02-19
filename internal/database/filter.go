package database

import "strings"

func toDatabaseFilter(inputFilter string) string {
	filter := inputFilter
	if filter == "" {
		filter = "%"
	} else {
		wildcardedFilter := strings.ReplaceAll(filter, " ", "%")
		filter = "%" + wildcardedFilter + "%"
	}
	return filter
}
