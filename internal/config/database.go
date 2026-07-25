package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
)

/* LoadEnv loads environment variables from a .env file */
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

/* ConnectDB establishes SQL and GORM connections */
func ConnectDB() {
	Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

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
