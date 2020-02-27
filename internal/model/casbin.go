package model

import (
	"github.com/google/uuid"
	"strings"
)

type CasbinRule struct {
	ID    uuid.UUID
	Type  string
	Value string
}

func (r CasbinRule) String() string {
	const prefixLine = ", "
	var sb strings.Builder

	sb.WriteString(r.Type)
	if len(r.Value) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(r.Value)
	}
	return sb.String()
}
