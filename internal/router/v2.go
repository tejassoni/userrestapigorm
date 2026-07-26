package router

import (
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func registerV2Routes(route *mux.Router, userHandler *handlers.UserHandler) {
	v2 := route.PathPrefix(config.API_PREFIX + "/" + config.API_VERSION).Subrouter()

	registerUserV2Routes(v2, userHandler)
	// registerAuthV2Routes(v2)
}
