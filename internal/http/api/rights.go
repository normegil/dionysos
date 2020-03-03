package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	httperror "github.com/normegil/dionysos/internal/http/error"
	middlewaresecurity "github.com/normegil/dionysos/internal/http/middleware/security"
	"github.com/normegil/dionysos/internal/security"
	"net/http"
)

type RightsController struct {
	DAO        RightsDAO
	ErrHandler httperror.HTTPErrorHandler
}

func (c RightsController) Route() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/current", c.current)
	return rt
}

func (c RightsController) current(w http.ResponseWriter, r *http.Request) {
	current := r.Context().Value(middlewaresecurity.KeyUser)
	if nil == current {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user := current.(security.User)
	rights, err := c.DAO.LoadForUser(user)
	if err != nil {
		c.ErrHandler.Handle(w, err)
		return
	}

	rightsBytes, err := json.Marshal(rights)
	if err != nil {
		c.ErrHandler.Handle(w, err)
		return
	}

	if _, err := w.Write(rightsBytes); nil != err {
		c.ErrHandler.Handle(w, err)
		return
	}
}

type RightsDAO interface {
	LoadForUser(security.User) ([]security.ResourceRights, error)
}
