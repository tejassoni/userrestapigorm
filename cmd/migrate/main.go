package main

import (
	"log"

	"userrestapigorm/internal/config"
	"userrestapigorm/internal/models"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("migration completed")
}
