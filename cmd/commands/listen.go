package commands

import (
	"context"
	"fmt"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/configuration"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	error2 "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func listen() (*cobra.Command, error) {
	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Launch dionysos server",
		Long:  `Launch dionysos server`,
		Run:   listenRun,
	}

	addressKey := configuration.KeyAddress
	listenCmd.Flags().StringP(addressKey.CommandLine.Name, addressKey.CommandLine.Shorthand, "0.0.0.0", addressKey.Description)
	if err := viper.BindPFlag(addressKey.Name, listenCmd.Flags().Lookup(addressKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", addressKey.Name, err)
	}

	portKey := configuration.KeyPort
	listenCmd.Flags().IntP(portKey.CommandLine.Name, portKey.CommandLine.Shorthand, 8080, portKey.Description)
	if err := viper.BindPFlag(portKey.Name, listenCmd.Flags().Lookup(portKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", portKey.Name, err)
	}

	return listenCmd, nil
}

func listenRun(_ *cobra.Command, _ []string) {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	addr := net.TCPAddr{
		IP:   net.ParseIP(viper.GetString(configuration.KeyAddress.Name)),
		Port: viper.GetInt(configuration.KeyPort.Name),
		Zone: "",
	}
	router := internalHTTP.NewRouter(newRoutes())
	handler := middleware.RequestLogger{Handler: router}
	closeHttpServer := internalHTTP.ListenAndServe(addr, handler)
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
	routes["/api"] = api.Controller{ErrHandler: error2.HTTPErrorHandler{}}.Routes()
	routes["/"] = http.FileServer(pkger.Dir("/website"))
	return routes
}
