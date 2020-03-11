package test

import (
	"encoding/json"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func FromJSONBody(t testing.TB, response []byte, v interface{}) {
	if err := json.Unmarshal(response, v); nil != err {
		t.Fatal(err)
	}
}

func ReadResponse(t testing.TB, response *httptest.ResponseRecorder) []byte {
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return respBytes
}

func HandlerErrorResponse(t testing.TB, code int, response []byte) {
	if code%200 > 99 {
		var respErr httperror.ErrorResponse
		FromJSONBody(t, response, respErr)
		t.Fatal(fmt.Errorf("received error response: %+v", respErr))
	}
}
