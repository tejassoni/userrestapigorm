package config

import (
	"fmt"
	"log"
	"regexp"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

// ConnectDB establishes the GORM connection.
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

	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting GORM:", err)
	}

	if err := configurePool(GormDB); err != nil {
		log.Fatal("Error configuring GORM connection pool:", err)
	}

	log.Println("Successfully connected to MySQL using GORM")
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

	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open MySQL server connection: %w", err)
	}

	adminSQLDB, err := adminDB.DB()
	if err != nil {
		return fmt.Errorf("access MySQL server connection: %w", err)
	}
	defer adminSQLDB.Close()

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci",
		name,
	)
	if err := adminDB.Exec(query).Error; err != nil {
		return fmt.Errorf("create database %q: %w", name, err)
	}

	return nil
}

func configurePool(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(DB_CONN_MAX_LIFETIME)
	sqlDB.SetConnMaxIdleTime(DB_CONN_MAX_IDLE_TIME)
	sqlDB.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
	sqlDB.SetMaxIdleConns(DB_MAX_IDLE_CONNS)

	return nil
}
