//nolint:funlen // There is no sense in limiting the size of a test function
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

func TestHTTPErrorHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		input       error
		expectation httperror.ErrorResponse
	}{
		{
			name:  "Simple error",
			input: errors.New("test error"),
			expectation: httperror.ErrorResponse{
				Code:   50000,
				Status: http.StatusInternalServerError,
				Error:  "test error",
			},
		},
		{
			name:  "Imbricated error",
			input: fmt.Errorf("error: %w", errors.New("test error")),
			expectation: httperror.ErrorResponse{
				Code:   50000,
				Status: http.StatusInternalServerError,
				Error:  "error: test error",
			},
		},
		{
			name: "HTTP Error",
			input: httperror.HTTPError{
				Code:   1,
				Status: 2,
				Err:    errors.New("http error"),
			},
			expectation: httperror.ErrorResponse{
				Code:   1,
				Status: 2,
				Error:  "http error",
			},
		},
		{
			name: "Imbricated HTTP Error",
			input: fmt.Errorf("error: %w", httperror.HTTPError{
				Code:   3,
				Status: 4,
				Err:    errors.New("http error"),
			}),
			expectation: httperror.ErrorResponse{
				Code:   3,
				Status: 4,
				Error:  "error: http error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := httperror.HTTPErrorHandler{}
			handler.Handle(recorder, test.input)

			result := recorder.Result()
			if test.expectation.Status != result.StatusCode {
				t.Errorf("Not expected status [header] {expected:%d;got:%d}", test.expectation.Status, result.StatusCode)
			}

			bytes, err := ioutil.ReadAll(result.Body)
			if nil != err {
				t.Fatal("could not read response body")
			}

			var errResponse httperror.ErrorResponse
			if err := json.Unmarshal(bytes, &errResponse); nil != err {
				t.Fatal("could not unmarshal response")
			}

			if test.expectation.Code != errResponse.Code {
				t.Errorf("Not expected response code {expected:%d;got:%d}", test.expectation.Code, errResponse.Code)
			}

			if test.expectation.Status != errResponse.Status {
				t.Errorf("Not expected response status {expected:%d;got:%d}", test.expectation.Status, errResponse.Status)
			}

			if test.expectation.Error != errResponse.Error {
				t.Errorf("Not expected response error message {expected:%s;got:%s}", test.expectation.Error, errResponse.Error)
			}
		})
	}
}
