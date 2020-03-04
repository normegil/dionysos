package security

import (
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

func AllHashAlgorithms() HashAlgorithms {
	return []HashAlgorithm{
		HashAlgorithmBcrypt14(),
	}
}
