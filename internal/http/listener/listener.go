package listener

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/casbin/casbin"
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"github.com/normegil/dionysos/internal/dao/database"
	"github.com/normegil/dionysos/internal/dao/database/versions"
	internalHTTP "github.com/normegil/dionysos/internal/http"
	"github.com/normegil/dionysos/internal/http/api"
	httperror "github.com/normegil/dionysos/internal/http/error"
	"github.com/normegil/dionysos/internal/http/middleware"
	securitymiddleware "github.com/normegil/dionysos/internal/http/middleware/security"
	"github.com/normegil/dionysos/internal/security"
	"github.com/normegil/dionysos/internal/security/authorization"
	"github.com/normegil/dionysos/internal/tools"
	"github.com/normegil/postgres"
	"io"
	"net/http"
)

type Configuration struct {
	APILogErrors bool
	Database     postgres.Configuration
}

type Listener struct {
	toClose       []io.Closer
	Configuration Configuration
}

func NewListener(cfg Configuration) *Listener {
	return &Listener{
		toClose:       make([]io.Closer, 0),
		Configuration: cfg,
	}
}

func (l Listener) Close() error {
	for _, closer := range l.toClose {
		tools.Close(closer)
	}
	return nil
}

func (l Listener) Load() (http.Handler, error) {
	db, err := l.initDatabase()
	if err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}
	l.toClose = append(l.toClose, db)
	return l.loadServerHandler(db), nil
}

func (l Listener) initDatabase() (*sql.DB, error) {
	dbCfg := l.Configuration.Database
	db, err := postgres.New(postgres.Configuration{
		Address:            dbCfg.Address,
		Port:               dbCfg.Port,
		User:               dbCfg.User,
		Password:           dbCfg.Password,
		Database:           dbCfg.Database,
		RequiredExtentions: append(l.Configuration.Database.RequiredExtentions, database.Extentions()...),
	})
	if err != nil {
		return nil, fmt.Errorf("creating database connection failed: %w", err)
	}

	manager, err := versions.NewSyncer(db)
	if err != nil {
		tools.Close(db)
		return db, fmt.Errorf("instantiate version manager: %w", err)
	}
	if err = manager.UpgradeAll(); nil != err {
		tools.Close(db)
		return db, fmt.Errorf("upgrading database: %w", err)
	}
	return db, nil
}

func (l Listener) loadServerHandler(db *sql.DB) http.Handler {
	sessionManager := scs.New()
	router := internalHTTP.NewRouter(l.route(db))
	sessionHandler := securitymiddleware.SessionHandler{
		SessionManager:       sessionManager,
		RequestAuthenticator: newRequestAuthenticator(database.UserDAO{DB: db}, sessionManager),
		ErrHandler:           httperror.HTTPErrorHandler{LogUserError: l.Configuration.APILogErrors},
		UserDAO:              database.UserDAO{DB: db},
		Handler:              router,
	}
	anonymousUserSetter := securitymiddleware.AnonymousUserSetter{Handler: sessionHandler}
	handler := middleware.RequestLogger{Handler: anonymousUserSetter}
	return handler
}

func (l Listener) route(db *sql.DB) map[string]http.Handler {
	itemDAO := &database.ItemDAO{DB: db}
	storageDAO := &database.StorageDAO{DB: db}
	casbinDAO := &database.CasbinDAO{
		DB:      db,
		RoleDAO: &security.NilRoleDAO{RoleDAO: &database.RoleDAO{DB: db}},
	}

	authorizer := newAuthorizer(casbinDAO)
	searchDAO := database.SearchDAO{
		Searchables: []database.Searcheable{
			itemDAO,
			storageDAO,
		},
		Authorizer: authorizer,
	}

	errorHandler := httperror.HTTPErrorHandler{LogUserError: l.Configuration.APILogErrors}
	authorizationHandler := securitymiddleware.AuthorizationHandler{
		Authorizer: authorizer,
		ErrHandler: errorHandler,
	}
	apiRoutes := make(map[string]http.Handler)
	apiRoutes["/storages"] = api.StorageController{
		StorageDAO: storageDAO,
		ErrHandler: errorHandler,
	}.Route(authorizationHandler)
	apiRoutes["/items"] = api.ItemController{
		ItemDAO:    itemDAO,
		ErrHandler: errorHandler,
	}.Route(authorizationHandler)
	apiRoutes["/searches"] = api.SearchController{
		DAO:        searchDAO,
		ErrHandler: errorHandler,
	}.Route()
	userCtrl := api.UserController{ErrHandler: errorHandler}
	apiRoutes["/users"] = userCtrl.Route()
	rightsCtrl := api.RightsController{
		DAO:        casbinDAO,
		ErrHandler: errorHandler,
	}
	apiRoutes["/rights"] = rightsCtrl.Route()
	apiCtrl := internalHTTP.MultiController{
		Routes: apiRoutes,
		OnRegister: func(rt *chi.Mux) {
			rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				http.StripPrefix("/api", http.FileServer(pkger.Dir("/api"))).ServeHTTP(w, r)
			})
		},
	}

	routes := make(map[string]http.Handler)
	routes["/api"] = apiCtrl.Route()
	routes["/"] = http.FileServer(pkger.Dir("/website/dist"))
	routes["/auth"] = api.AuthController{
		ErrHandler:     errorHandler,
		UserController: userCtrl,
	}.Route()
	return routes
}

func newAuthorizer(dao authorization.PolicyDAO) authorization.CasbinAuthorizer {
	adapter := &authorization.Adapter{DAO: dao}
	enforcer := casbin.NewEnforcer(authorization.Model(), adapter)
	authorizer := authorization.CasbinAuthorizer{Enforcer: enforcer}
	return authorizer
}

func newRequestAuthenticator(userDAO security.UserDAO, sessionManager *scs.SessionManager) securitymiddleware.RequestAuthenticator {
	authenticator := security.Authenticator{DAO: userDAO}
	updater := securitymiddleware.AuthenticatedUserSessionUpdater{
		SessionManager: sessionManager,
	}
	requestAuthenticator := securitymiddleware.RequestAuthenticator{
		Authenticator:   authenticator,
		OnAuthenticated: updater.RenewSessionOnAuthenticatedUser,
	}
	return requestAuthenticator
}
