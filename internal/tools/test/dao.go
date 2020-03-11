package test

import "testing"

func Clear(t testing.TB, clearers []Clearer) {
	for _, clearer := range clearers {
		if err := clearer.Clear(); nil != err {
			t.Error(err)
		}
	}
}

type Clearer interface {
	Clear() error
}
