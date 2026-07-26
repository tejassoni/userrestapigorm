package main

import (
	"log"

	"userrestapigorm/internal/config"
	"userrestapigorm/internal/models"
)

func main() {
	config.ConnectDB()

	if err := config.GormDB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("migration completed")
}
