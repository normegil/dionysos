package test

import (
	"encoding/json"
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
