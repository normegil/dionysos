package main

import (
	"context"
	"github.com/markbates/pkger"
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
	internalHTTP.NewRouter(newRoutes())

	router := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api/", http.FileServer(pkger.Dir("/api"))))

	closeHttpServer := internalHTTP.ListenAndServe(addr, router)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
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
