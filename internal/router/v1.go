package router

import (
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func registerV1Routes(route *mux.Router, userHandler *handlers.UserHandler) {
	v1 := route.PathPrefix(config.API_PREFIX + "/" + config.API_VERSION).Subrouter()

	registerUserV1Routes(v1, userHandler)
	registerAuthV1Routes(v1)
}
