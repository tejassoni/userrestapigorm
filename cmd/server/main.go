package main

import (
	"log"
	"net/http"
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/logger"
	"userrestapigorm/internal/router"
	"userrestapigorm/internal/validator"
)

func main() {

	// logs
	logger.New()
	// Initialize the validator
	validator.New()

	config.ConnectDB()         // Establish a connection to the database
	httpRouter := router.New() // Register application routes

	log.Println("Server started on " + config.APP_URL + ":" + config.APP_PORT) // Log the server start message

	if err := http.ListenAndServe(":"+config.APP_PORT, httpRouter); err != nil {
		log.Fatal(err)
	}
}
