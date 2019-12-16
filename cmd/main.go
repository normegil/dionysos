package main

import (
	"context"
	"github.com/normegil/dionysos/internal/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	srv := http.ListenAndServe()

	<-stopHTTPServer

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); nil != err {
		panic(err)
	}
}
