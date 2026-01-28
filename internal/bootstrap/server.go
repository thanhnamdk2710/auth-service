package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
	router "github.com/thanhnamdk2710/auth-service/internal/presentation/http"
)

type ServerOptions struct {
	Config   *config.ServerConfig
	Logger   *logger.Logger
	Metrics  *metrics.Metrics
	Handlers *Handlers
}

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func NewServer(opts ServerOptions) *Server {
	routerDeps := router.RouterDeps{
		Logger:      opts.Logger,
		Metrics:     opts.Metrics,
		AuthHandler: opts.Handlers.Auth,
	}

	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + opts.Config.Port,
			Handler:        router.New(routerDeps),
			ReadTimeout:    opts.Config.ReadTimeout,
			WriteTimeout:   opts.Config.WriteTimeout,
			IdleTimeout:    opts.Config.IdleTimeout,
			MaxHeaderBytes: opts.Config.MaxHeaderBytes,
		},
		logger: opts.Logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Server starting", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
