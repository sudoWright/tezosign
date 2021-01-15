package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"msig/common/log"
	"msig/conf"
	"msig/models"
	"msig/services/rpc_client"
	"net/http"
	"time"
)

const (
	gracefulTimeout  = time.Second * 10
	actionsAPIPrefix = ""
)

type (
	API struct {
		router       *mux.Router
		server       *http.Server
		cfg          conf.Config
		provider     NetworkContextProvider
		queryDecoder *schema.Decoder
	}

	// Route stores an API route data
	Route struct {
		Path       string
		Method     string
		Func       func(http.ResponseWriter, *http.Request)
		Middleware []negroni.HandlerFunc
	}
)

type NetworkContextProvider interface {
	GetDb(models.Network) (*gorm.DB, error)
	GetRPCClient(net models.Network) (*rpc_client.Tezos, error)
}

func NewAPI(cfg conf.Config, provider NetworkContextProvider) *API {
	queryDecoder := schema.NewDecoder()
	queryDecoder.IgnoreUnknownKeys(true)
	api := &API{
		cfg:          cfg,
		provider:     provider,
		queryDecoder: queryDecoder,
	}
	api.initialize()
	return api
}

// Run starts the http server and binds the handlers.
func (api *API) Run() error {
	return api.startServe()
}

func (api *API) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), gracefulTimeout)
	return api.server.Shutdown(ctx)
}

func (api *API) Title() string {
	return "API"
}

func (api *API) initialize(handlerArr ...negroni.Handler) {
	api.router = mux.NewRouter().UseEncodedPath()

	wrapper := negroni.New()

	for _, handler := range handlerArr {
		wrapper.Use(handler)
	}

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "User-Env"},
	}))

	//Static file
	w := wrapper.With(negroni.Wrap(http.StripPrefix("/static", http.FileServer(http.Dir("./resources")))))
	api.router.PathPrefix("/static").Handler(w).Methods(http.MethodGet, http.MethodOptions)

	//public routes
	HandleActions(api.router, wrapper, actionsAPIPrefix, []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},

		{Path: "/{network}/contract/storage/init", Method: http.MethodPost, Func: api.ContractStorageInit},
		{Path: "/{network}/contract/storage/update", Method: http.MethodPost, Func: api.ContractStorageUpdate},
		{Path: "/{network}/contract/operation", Method: http.MethodPost, Func: api.ContractOperation},
		{Path: "/{network}/contract/operation/signature", Method: http.MethodPost, Func: api.ContractOperationSignature},
	})

	api.server = &http.Server{Addr: fmt.Sprintf(":%d", api.cfg.API.ListenOnPort), Handler: api.router}
}

func (api *API) startServe() error {
	log.Info("Start listening server on port", zap.Uint64("port", api.cfg.API.ListenOnPort))
	err := api.server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Warn("API server was closed")
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot run API service: %s", err.Error())
	}
	return nil
}

// HandleActions is used to handle all given routes
func HandleActions(router *mux.Router, wrapper *negroni.Negroni, prefix string, routes []*Route) {
	for _, r := range routes {
		w := wrapper.With()
		for _, m := range r.Middleware {
			w.Use(m)
		}

		w.Use(negroni.Wrap(http.HandlerFunc(r.Func)))
		router.Handle(prefix+r.Path, w).Methods(r.Method, "OPTIONS")
	}
}
