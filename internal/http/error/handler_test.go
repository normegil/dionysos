package error_test

import (
	"encoding/json"
	"errors"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//nolint:funlen // Test functions with arrays might be quite long
func TestHTTPErrorHandler_Handle(t *testing.T) {
	type expected struct {
		Code    int
		Status  int
		Message string
	}

	tests := []struct {
		name   string
		err    error
		expect expected
	}{
		{
			name: "Simple error",
			err:  errors.New("testerror"),
			expect: expected{
				Code:    50000,
				Status:  http.StatusInternalServerError,
				Message: "testerror",
			},
		},
		{
			name: "HTTP error",
			err: httperror.HTTPError{
				Status: 567,
				Code:   56789,
				Err:    errors.New("testerror"),
			},
			expect: expected{
				Status:  567,
				Code:    56789,
				Message: "testerror",
			},
		},
		{
			name: "Deep HTTP error",
			err: fmt.Errorf("wrapping: %w", httperror.HTTPError{
				Code:   45678,
				Status: 456,
				Err:    errors.New("testerror"),
			}),
			expect: expected{
				Status:  456,
				Code:    45678,
				Message: "wrapping: testerror",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			httperror.HTTPErrorHandler{}.Handle(w, test.err)
			respBodyBytes, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			var respError httperror.ErrorResponse
			if err = json.Unmarshal(respBodyBytes, &respError); nil != err {
				t.Fatal(err)
			}

			if respError.Code != test.expect.Code {
				t.Errorf("Code is not equal to expected {expected:%d;got:%d}", test.expect.Code, respError.Code)
			}
			if respError.Status != test.expect.Status {
				t.Errorf("Status is not equal to expected {expected:%d;got:%d}", test.expect.Status, respError.Status)
			}
			if respError.Error != test.expect.Message {
				t.Errorf("Message is not equal to expected {expected:%s;got:%s}", test.expect.Message, respError.Error)
			}
		})
	}
}
