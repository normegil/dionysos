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
	"github.com/normegil/postgres"
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

	databaseAddressKey := configuration.KeyDatabaseAddress
	listenCmd.Flags().StringP(databaseAddressKey.CommandLine.Name, databaseAddressKey.CommandLine.Shorthand, "localhost", databaseAddressKey.Description)
	if err := viper.BindPFlag(databaseAddressKey.Name, listenCmd.Flags().Lookup(databaseAddressKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databaseAddressKey.Name, err)
	}

	databasePortKey := configuration.KeyDatabasePort
	listenCmd.Flags().IntP(databasePortKey.CommandLine.Name, databasePortKey.CommandLine.Shorthand, 5432, databasePortKey.Description)
	if err := viper.BindPFlag(databasePortKey.Name, listenCmd.Flags().Lookup(databasePortKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databasePortKey.Name, err)
	}

	databaseUserKey := configuration.KeyDatabaseUser
	listenCmd.Flags().IntP(databaseUserKey.CommandLine.Name, databaseUserKey.CommandLine.Shorthand, 5432, databaseUserKey.Description)
	if err := viper.BindPFlag(databaseUserKey.Name, listenCmd.Flags().Lookup(databaseUserKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databaseUserKey.Name, err)
	}

	databasePasswordKey := configuration.KeyDatabasePassword
	listenCmd.Flags().IntP(databasePasswordKey.CommandLine.Name, databasePasswordKey.CommandLine.Shorthand, 5432, databasePasswordKey.Description)
	if err := viper.BindPFlag(databasePasswordKey.Name, listenCmd.Flags().Lookup(databasePasswordKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databasePasswordKey.Name, err)
	}

	databaseNameKey := configuration.KeyDatabaseName
	listenCmd.Flags().IntP(databaseNameKey.CommandLine.Name, databaseNameKey.CommandLine.Shorthand, 5432, databaseNameKey.Description)
	if err := viper.BindPFlag(databaseNameKey.Name, listenCmd.Flags().Lookup(databaseNameKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databaseNameKey.Name, err)
	}

	return listenCmd, nil
}

func listenRun(_ *cobra.Command, _ []string) {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	db, err := postgres.New(postgres.Configuration{
		Address:  viper.GetString(configuration.KeyDatabaseAddress.Name),
		Port:     viper.GetInt(configuration.KeyDatabasePort.Name),
		User:     viper.GetString(configuration.KeyDatabaseUser.Name),
		Password: viper.GetString(configuration.KeyDatabasePassword.Name),
		Database: viper.GetString(configuration.KeyDatabaseName.Name),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("creating database connection failed")
	}

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
	routes["/"] = http.FileServer(pkger.Dir("/website/dist"))
	return routes
}
