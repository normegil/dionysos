package error

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type ReturningErrorHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HTTPErrorHandler struct {
	handler ReturningErrorHandler
}

func (h HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handler.ServeHTTP(w, r); nil != err {
		httpError := HTTPError{
			Code:   500,
			Status: 5000,
			Error:  err.Error(),
			Time:   time.Now(),
		}
		bytes, marshalErr := json.Marshal(httpError)
		if marshalErr != nil {
			log.Error().Err(marshalErr).Interface("HTTPError", httpError).Msg("Could not marshal HTTPError")
		}
		if _, writeErr := w.Write(bytes); nil != writeErr {
			log.Error().Err(writeErr).Interface("HTTPError", httpError).Msg("Could not write response with HTTPError")
		}
	}
}
