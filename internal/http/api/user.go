package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	httperror "github.com/normegil/dionysos/internal/http/error"
	middlewaresecurity "github.com/normegil/dionysos/internal/http/middleware/security"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
)

type UserController struct {
	ErrHandler httperror.HTTPErrorHandler
}

func (c UserController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/current", c.current)
	return rt
}

func (c UserController) current(w http.ResponseWriter, r *http.Request) {
	// Authentication handled by middlewares
	current := r.Context().Value(middlewaresecurity.KeyUser)
	if nil == current {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user := current.(security.User)
	jsonUser, err := json.Marshal(user)
	if err != nil {
		c.ErrHandler.Handle(w, fmt.Errorf("marhsal current user %+v: %w", user, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(jsonUser); nil != err {
		c.ErrHandler.Handle(w, fmt.Errorf("write current user (%+v) response: %w", user, err))
		return
	}
}
