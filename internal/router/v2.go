package router

import (
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func registerV2Routes(route *mux.Router, userHandler *handlers.UserHandler, cfg *config.Config) {
	v2 := route.PathPrefix(cfg.API.Prefix + "/" + cfg.API.Version).Subrouter()

	registerUserV2Routes(v2, userHandler)
	// registerAuthV2Routes(v2)
}
