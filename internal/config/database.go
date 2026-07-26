package config

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

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

	if err := ensureDatabase(dbUser, dbPassword, dbHost, dbPort, dbName); err != nil {
		log.Fatal("Error creating database:", err)
	}

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

	configurePool(DB)

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

	configurePool(sqlDB)

	log.Println("Successfully connected to MySQL using database/sql and GORM")
}

// ensureDatabase creates DB_NAME when it does not already exist. It connects
// without selecting a database because MySQL cannot connect to a missing one.
func ensureDatabase(user, password, host, port, name string) error {
	if !regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString(name) {
		return fmt.Errorf("invalid DB_NAME %q: use only letters, numbers, and underscores", name)
	}

	adminDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
	)

	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		return fmt.Errorf("open MySQL server connection: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return fmt.Errorf("ping MySQL server: %w", err)
	}

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci",
		name,
	)
	if _, err := adminDB.Exec(query); err != nil {
		return fmt.Errorf("create database %q: %w", name, err)
	}

	return nil
}

func configurePool(db *sql.DB) {
	db.SetConnMaxLifetime(DB_CONN_MAX_LIFETIME)
	db.SetConnMaxIdleTime(DB_CONN_MAX_IDLE_TIME)
	db.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
	db.SetMaxIdleConns(DB_MAX_IDLE_CONNS)
}
