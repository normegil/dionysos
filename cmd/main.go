package main

import (
	"context"
	"github.com/normegil/dionysos/internal/http"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	srv := http.ListenAndServe(net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8080,
		Zone: "",
	}, http.NewRouter())

	<-stopHTTPServer

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); nil != err {
		panic(err)
	}
}
