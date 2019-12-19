package main

import (
	"context"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/configuration"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	logCfg "github.com/normegil/dionysos/internal/log"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg, err := configuration.NewConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("loading configuration")
	}

	logCfg.Configure()

	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	addr := net.TCPAddr{
		IP:   net.ParseIP(cfg.GetString(configuration.KeyAddress.String())),
		Port: cfg.GetInt(configuration.KeyPort.String()),
		Zone: "",
	}
	rt := internalHTTP.NewRouter(newRoutes())
	closeHttpServer := internalHTTP.ListenAndServe(addr, rt)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := closeHttpServer(ctx); nil != err {
			log.Fatal().Err(err).Msg("closing server failed")
		}
	}()

	<-stopHTTPServer
}

func newRoutes() map[string]http.Handler {
	routes := make(map[string]http.Handler)
	routes["/api"] = api.NewRouter()
	routes["/"] = http.FileServer(pkger.Dir("/website"))
	return routes
}
