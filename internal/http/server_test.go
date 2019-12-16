package http_test

import (
	"context"
	"fmt"
	"github.com/normegil/connectionutils"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/interval"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	listeningIP := net.ParseIP("127.0.0.1")
	tcpAddr := connectionutils.SelectPort(listeningIP, *interval.MustParseIntervalInteger("[18880;18890]"))

	srv := internalHTTP.ListenAndServe(tcpAddr)
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		if err := srv.Shutdown(ctx); nil != err {
			panic(err)
		}
	}()

	resp, err := http.Get(fmt.Sprintf("http://%s:%d", tcpAddr.IP.String(), tcpAddr.Port))
	if nil != err {
		t.Fatal(err)
	}
	expected := http.StatusNotFound
	if expected != resp.StatusCode {
		t.Errorf("Wrong status {Expected:%d;Got:%d}", expected, resp.StatusCode)
	}
}
