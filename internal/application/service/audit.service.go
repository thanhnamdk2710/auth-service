package service

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type AuditService interface {
	Log(ctx context.Context, log *entity.AuditLog)
	Start()
	Stop()
}

type AuditServiceConfig struct {
	BufferSize   int
	WorkerCount  int
	BatchSize    int
	FlushTimeout time.Duration
}

func DefaultAuditServiceConfig() AuditServiceConfig {
	return AuditServiceConfig{
		BufferSize:   1000,
		WorkerCount:  4,
		BatchSize:    50,
		FlushTimeout: 5 * time.Second,
	}
}

type auditService struct {
	repo    repository.AuditRepository
	logger  *logger.Logger
	config  AuditServiceConfig
	logChan chan *entity.AuditLog
	wg      sync.WaitGroup
	stopCh  chan struct{}
}

func NewAuditService(repo repository.AuditRepository, logger *logger.Logger, cfg AuditServiceConfig) AuditService {
	return &auditService{
		repo:    repo,
		logger:  logger,
		config:  cfg,
		logChan: make(chan *entity.AuditLog, cfg.BufferSize),
		stopCh:  make(chan struct{}),
	}
}

func (s *auditService) Log(ctx context.Context, log *entity.AuditLog) {
	select {
	case s.logChan <- log:
	default:
		s.logger.Warn("Audit log buffer full, dropping log",
			zap.String("action", string(log.Action)),
			zap.String("correlation_id", log.CorrelationID),
		)
	}
}

func (s *auditService) Start() {
	for i := 0; i < s.config.WorkerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
	s.logger.Info("Audit service started",
		zap.Int("workers", s.config.WorkerCount),
		zap.Int("buffer_size", s.config.BufferSize),
		zap.Int("batch_size", s.config.BatchSize),
	)
}

func (s *auditService) Stop() {
	close(s.stopCh)
	close(s.logChan)
	s.wg.Wait()
	s.logger.Info("Audit service stopped")
}

func (s *auditService) worker(id int) {
	defer s.wg.Done()

	batch := make([]*entity.AuditLog, 0, s.config.BatchSize)
	ticker := time.NewTicker(s.config.FlushTimeout)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, log := range batch {
			if err := s.repo.Create(ctx, log); err != nil {
				s.logger.Error("Failed to write audit log",
					zap.Error(err),
					zap.Int("worker_id", id),
					zap.String("action", string(log.Action)),
					zap.String("correlation_id", log.CorrelationID),
				)
			}
		}

		batch = batch[:0]
	}

	for {
		select {
		case log, ok := <-s.logChan:
			if !ok {
				flush()
				return
			}
			batch = append(batch, log)
			if len(batch) >= s.config.BatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-s.stopCh:
			for log := range s.logChan {
				batch = append(batch, log)
				if len(batch) >= s.config.BatchSize {
					flush()
				}
			}
			flush()
			return
		}
	}
}
