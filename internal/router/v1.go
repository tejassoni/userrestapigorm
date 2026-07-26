package router

import (
	"userrestapigorm/internal/config"

	"github.com/gorilla/mux"
)

func registerV1Routes(route *mux.Router) {
	v1 := route.PathPrefix(config.API_PREFIX + "/" + config.API_VERSION).Subrouter()

	registerUserV1Routes(v1)
	registerAuthV1Routes(v1)
}
