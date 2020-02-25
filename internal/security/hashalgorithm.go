package security

import (
	"fmt"
	"github.com/google/uuid"
)

type HashAlgorithm interface {
	ID() uuid.UUID
	HashAndSalt(password string) ([]byte, error)
	Validate(hash []byte, password string) error
}

type HashAlgorithms []HashAlgorithm

func (h HashAlgorithms) FindByID(id uuid.UUID) HashAlgorithm {
	for _, algorithm := range h {
		if id == algorithm.ID() {
			return algorithm
		}
	}
	return nil
}

var (
	AllHashAlgorithms = HashAlgorithms([]HashAlgorithm{
		HashAlgorithmBcrypt14,
	})
)

type invalidPasswordError struct {
	Original error
}

func (e invalidPasswordError) Error() string {
	return fmt.Errorf("invalid password: %w", e.Original).Error()
}
