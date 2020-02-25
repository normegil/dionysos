package model

import "github.com/google/uuid"

type HashAlgorithm struct {
	ID               uuid.UUID
	AlgorithmOptions interface{}
}

var (
	Bcrypt14 = HashAlgorithm{
		ID:               uuid.MustParse("01d1de6c-fa67-4caa-84da-684dc5640626"),
		AlgorithmOptions: BCryptOptions{Cost: 14},
	}
)

type BCryptOptions struct {
	Cost int
}
