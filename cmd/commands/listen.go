package commands

import (
	"context"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/configuration"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	RootCmd.AddCommand(listenCmd)

	addressKey := configuration.KeyAddress
	listenCmd.Flags().StringP(addressKey.CommandLine.Name, addressKey.CommandLine.Shorthand, "0.0.0.0", addressKey.Description)
	if err := viper.BindPFlag(addressKey.Name, listenCmd.Flags().Lookup(addressKey.CommandLine.Name)); err != nil {
		log.Fatal().Err(err).Str("paramName", addressKey.Name).Msg("Binding parameter")
	}

	portKey := configuration.KeyPort
	listenCmd.Flags().IntP(portKey.CommandLine.Name, portKey.CommandLine.Shorthand, 8080, portKey.Description)
	if err := viper.BindPFlag(portKey.Name, listenCmd.Flags().Lookup(portKey.CommandLine.Name)); err != nil {
		log.Fatal().Err(err).Str("paramName", portKey.Name).Msg("Binding parameter")
	}
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Launch dionysos server",
	Long:  `Launch dionysos server`,
	Run: func(cmd *cobra.Command, args []string) {
		stopHTTPServer := make(chan os.Signal, 1)
		signal.Notify(stopHTTPServer, os.Interrupt)

		addr := net.TCPAddr{
			IP:   net.ParseIP(viper.GetString(configuration.KeyAddress.Name)),
			Port: viper.GetInt(configuration.KeyPort.Name),
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
	},
}

func newRoutes() map[string]http.Handler {
	routes := make(map[string]http.Handler)
	routes["/api"] = api.NewRouter()
	routes["/"] = http.FileServer(pkger.Dir("/website"))
	return routes
}
