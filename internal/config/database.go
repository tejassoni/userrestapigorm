package config

import (
	"fmt"
	"regexp"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

// ConnectDB establishes the GORM connection.
func ConnectDB() (*gorm.DB, error) {
	cfg := Load()
	databaseConfig := cfg.Database

	dbUser := databaseConfig.Username
	dbPassword := databaseConfig.Password
	dbHost := databaseConfig.Host
	dbPort := strconv.Itoa(databaseConfig.Port)
	dbName := databaseConfig.Database

	if err := ensureDatabase(dbUser, dbPassword, dbHost, dbPort, dbName); err != nil {
		return nil, fmt.Errorf("create database: %w", err)
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
		return nil, fmt.Errorf("open GORM connection: %w", err)
	}

	if err := configurePool(GormDB, databaseConfig); err != nil {
		return nil, fmt.Errorf("configure connection pool: %w", err)
	}

	return GormDB, nil
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

func configurePool(gormDB *gorm.DB, config DatabaseConfig) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(config.MaxConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	return nil
}
