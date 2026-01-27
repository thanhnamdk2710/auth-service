package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/application/service"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
	router "github.com/thanhnamdk2710/auth-service/internal/presentation/http"
)

type ServerOptions struct {
	Port         string
	Logger       *logger.Logger
	Metrics      *metrics.Metrics
	DB           *Database
	AuditService service.AuditService
}

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func NewServer(opts ServerOptions) *Server {
	routerDeps := router.Deps{
		Logger:       opts.Logger,
		Metrics:      opts.Metrics,
		DB:           opts.DB.Conn(),
		AuditService: opts.AuditService,
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + opts.Port,
			Handler:      router.New(routerDeps),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: opts.Logger,
	}
}

func (s *Server) Start() {
	s.logger.Info("Server starting", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Server failed", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
