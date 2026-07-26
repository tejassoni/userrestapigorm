package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"
	"userrestapigorm/internal/logger"
	"userrestapigorm/internal/repository"
	"userrestapigorm/internal/router"
	"userrestapigorm/internal/service"
)

type Application struct {
	server  *http.Server
	logger  *slog.Logger // or your concrete logger type
	cleanup func()
}

func Build() (*Application, error) {

	appLogger := logger.New()

	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("sql.DB error: %w", err)
	}

	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("error closing sql.DB", "error", err)
		}
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, appLogger)
	httpRouter := router.New(userHandler)

	srv := &http.Server{
		Addr:              ":" + config.APP_PORT,
		Handler:           httpRouter,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Application{
		server:  srv,
		logger:  appLogger,
		cleanup: cleanup,
	}, nil
}

func (a *Application) Run() error {
	defer a.cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		a.logger.Info("server started", "port", config.APP_PORT)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("server listener error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	a.logger.Info("shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.server.Shutdown(shutdownCtx)
}
