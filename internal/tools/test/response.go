package test

import (
	"encoding/json"
	"fmt"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func FromJSONBody(t testing.TB, response *httptest.ResponseRecorder, v interface{}) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(bodyBytes, v); nil != err {
		t.Fatal(err)
	}
}

func HandlerErrorResponse(t testing.TB, response *httptest.ResponseRecorder) {
	if response.Code%200 > 99 {
		respBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		var respErr httperror.ErrorResponse
		if err := json.Unmarshal(respBytes, err); nil != err {
			t.Fatal(err)
		}
		t.Fatal(fmt.Errorf("received error response: %+v", respErr))
	}
}
