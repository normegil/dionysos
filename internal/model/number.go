package model

import "fmt"

type Natural struct {
	number int
}

func NewNatural(number int) (*Natural, error) {
	if number < 0 {
		return nil, fmt.Errorf("%d is not a natural number", number)
	}
	return &Natural{number: number}, nil
}

func MustNewNatural(number int) *Natural {
	natural, err := NewNatural(number)
	if err != nil {
		panic(err)
	}
	return natural
}

func (n Natural) Number() int {
	return n.number
}
