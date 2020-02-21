package dionysos

import (
	"github.com/google/uuid"
)

type Storage struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
