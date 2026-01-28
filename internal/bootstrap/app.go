package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
)

type App struct {
	cfg      *config.Config
	logger   *logger.Logger
	db       *Database
	metrics  *metrics.Metrics
	services *Services
	handlers *Handlers
	server   *Server
}

func New() (*App, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New(&logger.Config{
		Level:       cfg.Server.LogLevel,
		Environment: cfg.Server.Environment,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		logger: log,
	}, nil
}

func (a *App) Run() error {
	ctx := context.Background()
	a.logger.Info("Starting application...")

	if err := a.initDatabase(ctx); err != nil {
		return err
	}

	a.initMetrics()
	a.initServices()
	a.initHandlers()
	a.initServer()

	serverErr := make(chan error, 1)
	go func() {
		if err := a.server.Start(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		a.logger.Error("Server error", zap.Error(err))
		a.shutdown(ctx)
		return err
	case <-a.waitForSignal():
		a.shutdown(ctx)
		return nil
	}
}

func (a *App) initDatabase(ctx context.Context) error {
	db, err := NewDatabase(ctx, a.cfg.DB)
	if err != nil {
		a.logger.FatalCtx(ctx, "Failed to connect to database", zap.Error(err))
		return err
	}
	a.db = db
	a.logger.Info("Database connected")
	return nil
}

func (a *App) initMetrics() {
	a.metrics = metrics.New(prometheus.DefaultRegisterer)
	a.metrics.RegisterDBStats(prometheus.DefaultRegisterer, a.db.SQL())
}

func (a *App) initServices() {
	a.services = NewServices(a.db, a.logger)
	a.services.Start()
}

func (a *App) initHandlers() {
	a.handlers = NewHandlers(a.db, a.services.Audit(), a.logger)
}

func (a *App) initServer() {
	a.server = NewServer(ServerOptions{
		Config:   a.cfg.Server,
		Logger:   a.logger,
		Metrics:  a.metrics,
		Handlers: a.handlers,
	})
}

func (a *App) waitForSignal() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func (a *App) shutdown(ctx context.Context) {
	a.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.ErrorCtx(ctx, "Server forced to shutdown", zap.Error(err))
	}

	a.logger.Info("Stopping services...")
	a.services.Stop()

	a.logger.Info("Closing database connection...")
	if err := a.db.Close(); err != nil {
		a.logger.ErrorCtx(ctx, "Error closing database", zap.Error(err))
	}

	a.logger.Sync()
	a.logger.Info("Application stopped")
}
