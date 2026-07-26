package router

import (
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func registerV1Routes(route *mux.Router, userHandler *handlers.UserHandler, cfg *config.Config) {
	v1 := route.PathPrefix(cfg.API.Prefix + "/" + cfg.API.Version).Subrouter()

	registerUserV1Routes(v1, userHandler)
	registerAuthV1Routes(v1)
}
