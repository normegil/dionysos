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
	tcpAddr, closeHttpServer := newHTTPServer(t)
	defer closeHttpServer()

	resp, err := http.Get(fmt.Sprintf("http://%s:%d", tcpAddr.IP.String(), tcpAddr.Port))
	if nil != err {
		t.Fatal(err)
	}
	expected := http.StatusNotFound
	if expected != resp.StatusCode {
		t.Errorf("Wrong status {Expected:%d;Got:%d}", expected, resp.StatusCode)
	}
}

func newHTTPServer(t testing.TB) (net.TCPAddr, func()) {
	listeningIP := net.ParseIP("127.0.0.1")
	tcpAddr := connectionutils.SelectPort(listeningIP, *interval.MustParseIntervalInteger("[18870;18890]"))

	closeHTTPServer := internalHTTP.ListenAndServe(tcpAddr, http.DefaultServeMux)
	return tcpAddr, func() {
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		if err := closeHTTPServer(ctx); nil != err {
			t.Fatal(fmt.Errorf("close http server: %w", err))
		}
	}
}
