package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
)

/* ConnectDB establishes SQL and GORM connections */
func ConnectDB() {

	dbUser := DB_USER         // Get the database user from environment variables config.go
	dbPassword := DB_PASSWORD // Get the database password from environment variables config.go
	dbHost := DB_HOST         // Get the database host from environment variables config.go
	dbPort := DB_PORT         // Get the database port from environment variables config.go
	dbName := DB_NAME         // Get the database name from environment variables config.go

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// ---------------------------
	// database/sql connection
	// ---------------------------
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(2 * time.Minute)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	if err = DB.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// ---------------------------
	// GORM connection
	// ---------------------------
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting GORM:", err)
	}

	// Apply the same pool settings
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	log.Println("Successfully connected to MySQL using database/sql and GORM")
}
