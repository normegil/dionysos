package commands

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/configuration"
	"github.com/normegil/dionysos/internal/database"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware"
	securitymiddleware "github.com/normegil/dionysos/internal/http/middleware/security"
	"github.com/normegil/dionysos/internal/security"
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

	db, err := InitDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("initializing database")
	}
	defer closeDatabase(db)

	addr := net.TCPAddr{
		IP:   net.ParseIP(viper.GetString(configuration.KeyAddress.Name)),
		Port: viper.GetInt(configuration.KeyPort.Name),
		Zone: "",
	}
	if err != nil {
		log.Fatal().Err(err).Msg("creating routing rules")
	}
	closeHttpServer := internalHTTP.ListenAndServe(addr, ToServerHandler(db))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := closeHttpServer(ctx); nil != err {
			log.Fatal().Err(err).Msg("closing server failed")
		}
	}()

	<-stopHTTPServer
}

func ToServerHandler(db *sql.DB) middleware.RequestLogger {
	sessionManager := scs.New()
	router := internalHTTP.NewRouter(Route(db, sessionManager))
	sessionHandler := securitymiddleware.SessionHandler{
		SessionManager:       sessionManager,
		RequestAuthenticator: NewRequestAuthenticator(database.UserDAO{DB: db}, sessionManager),
		ErrHandler:           ErrorHandler(),
		Handler:              router,
	}
	handler := middleware.RequestLogger{Handler: sessionHandler}
	return handler
}

func closeDatabase(db *sql.DB) {
	if err := db.Close(); nil != err {
		log.Error().Err(err).Msg("Could not close database connection")
	}
}

func getDatabaseConfiguration() postgres.Configuration {
	return postgres.Configuration{
		Address:            viper.GetString(configuration.KeyDatabaseAddress.Name),
		Port:               viper.GetInt(configuration.KeyDatabasePort.Name),
		User:               viper.GetString(configuration.KeyDatabaseUser.Name),
		Password:           viper.GetString(configuration.KeyDatabasePassword.Name),
		Database:           viper.GetString(configuration.KeyDatabaseName.Name),
		RequiredExtentions: database.Extentions,
	}
}

func InitDatabase() (*sql.DB, error) {
	dbCfg := getDatabaseConfiguration()
	db, err := postgres.New(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("creating database connection failed: %w", err)
	}

	manager, err := database.NewVersionManager(db)
	if err != nil {
		closeDatabase(db)
		return db, fmt.Errorf("instantiate version manager: %w", err)
	}
	if err = manager.UpgradeAll(); nil != err {
		closeDatabase(db)
		return db, fmt.Errorf("upgrading database: %w", err)
	}
	if err := database.InsertDummyData(db); nil != err {
		closeDatabase(db)
		return db, fmt.Errorf("inserting dummy data: %w", err)
	}
	return db, nil
}

func ErrorHandler() httperror.HTTPErrorHandler {
	return httperror.HTTPErrorHandler{LogUserError: viper.GetBool(configuration.KeyAPIShowError.Name)}
}

func NewRequestAuthenticator(userDAO security.UserDAO, sessionManager *scs.SessionManager) securitymiddleware.RequestAuthenticator {
	authenticator := security.DatabaseAuthentication{DAO: userDAO}
	updater := securitymiddleware.AuthenticatedUserSessionUpdater{
		SessionManager: sessionManager,
	}
	requestAuthenticator := securitymiddleware.RequestAuthenticator{
		Authenticator:   authenticator,
		OnAuthenticated: updater.RenewSessionOnAuthenticatedUser,
	}
	return requestAuthenticator
}

func Route(db *sql.DB, sessionManager *scs.SessionManager) map[string]http.Handler {
	itemDAO := &database.ItemDAO{DB: db}
	storageDAO := &database.StorageDAO{DB: db}
	searchDAO := database.SearchDAO{Searchables: []database.Searcheable{
		itemDAO,
		storageDAO,
	}}

	apiRoutes := make(map[string]http.Handler)
	errorHandler := ErrorHandler()
	apiRoutes["/storages"] = api.StorageController{
		StorageDAO: storageDAO,
		ErrHandler: errorHandler,
	}.Route()
	apiRoutes["/items"] = api.ItemController{
		ItemDAO:    itemDAO,
		ErrHandler: errorHandler,
	}.Route()
	apiRoutes["/searches"] = api.SearchController{
		DAO:        searchDAO,
		ErrHandler: errorHandler,
	}.Route()
	apiCtrl := internalHTTP.MultiController{
		Routes: apiRoutes,
		OnRegister: func(rt *chi.Mux) {
			rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				http.StripPrefix("/api", http.FileServer(pkger.Dir("/api"))).ServeHTTP(w, r)
			})
		},
	}

	userDAO := &database.UserDAO{DB: db}

	routes := make(map[string]http.Handler)
	routes["/api"] = apiCtrl.Route()
	routes["/"] = http.FileServer(pkger.Dir("/website/dist"))
	routes["/auth"] = api.AuthController{
		ErrHandler:           errorHandler,
		RequestAuthenticator: NewRequestAuthenticator(userDAO, sessionManager),
	}.Route()
	return routes
}
