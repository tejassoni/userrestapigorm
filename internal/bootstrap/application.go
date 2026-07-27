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
	"userrestapigorm/internal/validator"
)

type Application struct {
	config  *config.Config
	server  *http.Server
	logger  *slog.Logger
	cleanup *Cleanup
}

func Build() (app *Application, err error) {

	cleanup := NewCleanup()

	// Temporary logger until config is loaded
	appLogger := logger.New()

	// Cleanup partially initialized resources if Build fails
	defer func() {
		if err != nil {
			appLogger.Warn("bootstrap failed, cleaning resources")
			cleanup.Run(appLogger)
		}
	}()

	cfg := config.Load()

	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	cleanup.Add("mysql", sqlDB.Close)

	// --------------------------
	// Modules
	// --------------------------

	userRepo := repository.NewUserRepository(db)

	userValidator := validator.New(userRepo)

	userService := service.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(
		userService,
		appLogger,
		userValidator,
	)

	httpRouter := router.New(userHandler, cfg)

	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           httpRouter,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	app = &Application{
		config:  cfg,
		server:  server,
		logger:  appLogger,
		cleanup: cleanup,
	}

	return
}

func (a *Application) Run() error {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)

	go func() {

		a.logger.Info(
			"HTTP server started",
			"port", a.config.Server.Port,
		)

		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			errCh <- err
		}
	}()

	select {

	case err := <-errCh:

		a.cleanup.Run(a.logger)

		return fmt.Errorf("http server: %w", err)

	case <-ctx.Done():

		a.logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {

			a.cleanup.Run(a.logger)

			return fmt.Errorf("shutdown: %w", err)
		}

		a.cleanup.Run(a.logger)

		a.logger.Info("application stopped")

		return nil
	}
}
