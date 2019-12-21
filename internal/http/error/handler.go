package error

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HTTPErrorHandler struct {
}

func (h HTTPErrorHandler) Handle(w http.ResponseWriter, err error) {
	httpError := &HTTPError{}
	if !errors.As(err, httpError) {
		httpError = &HTTPError{
			Code:   50000,
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	resp := httpError.toResponse()
	bytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Interface("HTTPError", resp).Msg("Could not marshal HTTPError")
	}
	if _, writeErr := w.Write(bytes); nil != writeErr {
		log.Error().Err(writeErr).Interface("HTTPError", resp).Msg("Could not write response with HTTPError")
	}
}
