package commands

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/configuration"
	"github.com/normegil/dionysos/internal/dao/database"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/listener"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"os"
	"os/signal"
	"time"
)

//nolint:funlen // Main function is quite long but highly repetitive with subtle differences. Limiting function by size isn't required as the coplexity is still low.
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

	apiErrorShowKey := configuration.KeyAPIShowError
	listenCmd.Flags().BoolP(apiErrorShowKey.CommandLine.Name, apiErrorShowKey.CommandLine.Shorthand, false, apiErrorShowKey.Description)
	if err := viper.BindPFlag(apiErrorShowKey.Name, listenCmd.Flags().Lookup(apiErrorShowKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", apiErrorShowKey.Name, err)
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
	listenCmd.Flags().StringP(databaseUserKey.CommandLine.Name, databaseUserKey.CommandLine.Shorthand, "postgres", databaseUserKey.Description)
	if err := viper.BindPFlag(databaseUserKey.Name, listenCmd.Flags().Lookup(databaseUserKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databaseUserKey.Name, err)
	}

	databasePasswordKey := configuration.KeyDatabasePassword
	listenCmd.Flags().StringP(databasePasswordKey.CommandLine.Name, databasePasswordKey.CommandLine.Shorthand, "postgres", databasePasswordKey.Description)
	if err := viper.BindPFlag(databasePasswordKey.Name, listenCmd.Flags().Lookup(databasePasswordKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databasePasswordKey.Name, err)
	}

	databaseNameKey := configuration.KeyDatabaseName
	listenCmd.Flags().StringP(databaseNameKey.CommandLine.Name, databaseNameKey.CommandLine.Shorthand, "dionysos", databaseNameKey.Description)
	if err := viper.BindPFlag(databaseNameKey.Name, listenCmd.Flags().Lookup(databaseNameKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", databaseNameKey.Name, err)
	}

	dummyDataKey := configuration.KeyDummyData
	listenCmd.Flags().BoolP(dummyDataKey.CommandLine.Name, dummyDataKey.CommandLine.Shorthand, false, dummyDataKey.Description)
	if err := viper.BindPFlag(dummyDataKey.Name, listenCmd.Flags().Lookup(dummyDataKey.CommandLine.Name)); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", dummyDataKey.Name, err)
	}

	return listenCmd, nil
}

func listenRun(_ *cobra.Command, _ []string) {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	dbCfg := getDatabaseConfiguration()
	handler, err := listener.NewListener(listener.Configuration{
		APILogErrors: viper.GetBool(configuration.KeyAPIShowError.Name),
		Database:     dbCfg,
	}).Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Initializing Listener")
	}

	if viper.GetBool(configuration.KeyDummyData.Name) {
		db, err := postgres.New(dbCfg)
		if err != nil {
			log.Fatal().Err(err).Msg("initializing connection to database to insert dummy data")
		}
		if err = insertDummyData(db); nil != err {
			log.Fatal().Err(err).Msg("inserting dummy data")
		}

	}

	addr := net.TCPAddr{
		IP:   net.ParseIP(viper.GetString(configuration.KeyAddress.Name)),
		Port: viper.GetInt(configuration.KeyPort.Name),
		Zone: "",
	}
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

func getDatabaseConfiguration() postgres.Configuration {
	return postgres.Configuration{
		Address:            viper.GetString(configuration.KeyDatabaseAddress.Name),
		Port:               viper.GetInt(configuration.KeyDatabasePort.Name),
		User:               viper.GetString(configuration.KeyDatabaseUser.Name),
		Password:           viper.GetString(configuration.KeyDatabasePassword.Name),
		Database:           viper.GetString(configuration.KeyDatabaseName.Name),
		RequiredExtentions: make([]string, 0),
	}
}

func insertDummyData(db *sql.DB) error {
	storageDAO := &database.StorageDAO{Querier: db}
	itemDAO := &database.ItemDAO{Querier: db}
	userDAO := &database.UserDAO{Querier: db}

	log.Info().Msg("insert dummy data")
	const dummyItemNb = 50

	rolename := "user"
	userrole, err := database.RoleDAO{Querier: db}.LoadByName(rolename)
	if err != nil {
		return fmt.Errorf("load '%s' role: %w", rolename, err)
	}

	user, err := security.NewUser("user", "user", *userrole)
	if err != nil {
		return fmt.Errorf("creating 'user' user: %w", err)
	}
	if err := userDAO.Insert(*user); nil != err {
		return fmt.Errorf("inserting %s: %w", user.Name, err)
	}

	for i := 0; i < dummyItemNb; i++ {
		storage := dionysos.Storage{Name: gofakeit.Company()}
		if err := storageDAO.Insert(storage); nil != err {
			return fmt.Errorf("inserting %+v: %w", storage, err)
		}
	}

	for i := 0; i < dummyItemNb; i++ {
		item := dionysos.Item{Name: gofakeit.BeerName()}
		if err := itemDAO.Insert(item); nil != err {
			return fmt.Errorf("inserting %+v: %w", item, err)
		}
	}
	return nil
}
