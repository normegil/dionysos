package main

import (
	"context"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8080,
		Zone: "",
	}
	rt := internalHTTP.NewRouter(newRoutes())
	closeHttpServer := internalHTTP.ListenAndServe(addr, rt)
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := closeHttpServer(ctx); nil != err {
			panic(err)
		}
	}()

	<-stopHTTPServer
}

func newRoutes() map[string]http.Handler {
	routes := make(map[string]http.Handler)
	routes["/api"] = api.NewRouter()
	return routes
}
