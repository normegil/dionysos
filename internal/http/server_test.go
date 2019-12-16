package http_test

import (
	"context"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	srv := internalHTTP.ListenAndServe()
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		if err := srv.Shutdown(ctx); nil != err {
			panic(err)
		}
	}()

	resp, err := http.Get("http://localhost:8080")
	if nil != err {
		t.Fatal(err)
	}
	expected := http.StatusNotFound
	if expected != resp.StatusCode {
		t.Errorf("Wrong status {Expected:%d;Got:%d}", expected, resp.StatusCode)
	}
}
