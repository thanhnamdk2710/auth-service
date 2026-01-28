package audit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type Config struct {
	BufferSize      int
	WorkerCount     int
	BatchSize       int
	FlushTimeout    time.Duration
	ShutdownTimeout time.Duration
}

const (
	DefaultBufferSize       = 1000
	DefaultWorkerCount      = 4
	DefaultBatchSize        = 50
	DefaultFlushTimeoutSec  = 5
	DefaultShutdownTimeout  = 30 * time.Second
)

func DefaultConfig() Config {
	return Config{
		BufferSize:      DefaultBufferSize,
		WorkerCount:    DefaultWorkerCount,
		BatchSize:       DefaultBatchSize,
		FlushTimeout:    DefaultFlushTimeoutSec * time.Second,
		ShutdownTimeout: DefaultShutdownTimeout,
	}
}

type AsyncLogger struct {
	repo           repository.AuditRepository
	log            *logger.Logger
	config         Config
	logChan        chan *entity.AuditLog
	wg             sync.WaitGroup
	stopCh         chan struct{}
	stopped        atomic.Bool
	droppedCounter atomic.Int64
}

func NewAsyncLogger(repo repository.AuditRepository, log *logger.Logger, cfg Config) *AsyncLogger {
	return &AsyncLogger{
		repo:    repo,
		log:     log,
		config:  cfg,
		logChan: make(chan *entity.AuditLog, cfg.BufferSize),
		stopCh:  make(chan struct{}),
	}
}

func (a *AsyncLogger) Log(ctx context.Context, auditLog *entity.AuditLog) {
	if a.stopped.Load() {
		a.log.Warn("Audit logger stopped, dropping log",
			zap.String("action", string(auditLog.Action)),
			zap.String("correlation_id", auditLog.CorrelationID),
		)
		return
	}

	select {
	case a.logChan <- auditLog:
	case <-ctx.Done():
		a.log.Warn("Context cancelled, dropping audit log",
			zap.String("action", string(auditLog.Action)),
			zap.String("correlation_id", auditLog.CorrelationID),
		)
	default:
		a.droppedCounter.Add(1)
		a.log.Warn("Audit log buffer full, dropping log",
			zap.String("action", string(auditLog.Action)),
			zap.String("correlation_id", auditLog.CorrelationID),
			zap.Int64("total_dropped", a.droppedCounter.Load()),
		)
	}
}

func (a *AsyncLogger) Start() {
	for i := 0; i < a.config.WorkerCount; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}
	a.log.Info("Audit logger started",
		zap.Int("workers", a.config.WorkerCount),
		zap.Int("buffer_size", a.config.BufferSize),
		zap.Int("batch_size", a.config.BatchSize),
	)
}

func (a *AsyncLogger) Stop() {
	a.stopped.Store(true)
	close(a.stopCh)

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.log.Info("Audit logger stopped gracefully",
			zap.Int64("total_dropped", a.droppedCounter.Load()),
		)
	case <-time.After(a.config.ShutdownTimeout):
		a.log.Warn("Audit logger shutdown timed out",
			zap.Duration("timeout", a.config.ShutdownTimeout),
			zap.Int64("total_dropped", a.droppedCounter.Load()),
		)
	}
}

func (a *AsyncLogger) DroppedCount() int64 {
	return a.droppedCounter.Load()
}

func (a *AsyncLogger) QueueSize() int {
	return len(a.logChan)
}

func (a *AsyncLogger) worker(id int) {
	defer a.wg.Done()

	batch := make([]*entity.AuditLog, 0, a.config.BatchSize)
	ticker := time.NewTicker(a.config.FlushTimeout)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.repo.CreateBatch(ctx, batch); err != nil {
			a.log.Error("Failed to write audit log batch",
				zap.Error(err),
				zap.Int("worker_id", id),
				zap.Int("batch_size", len(batch)),
			)
		}

		batch = batch[:0]
	}

	for {
		select {
		case auditLog := <-a.logChan:
			if auditLog == nil {
				flush()
				return
			}
			batch = append(batch, auditLog)
			if len(batch) >= a.config.BatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-a.stopCh:
			a.drainWithTimeout(id, &batch, flush)
			return
		}
	}
}

func (a *AsyncLogger) drainWithTimeout(workerID int, batch *[]*entity.AuditLog, flush func()) {
	drainTimeout := time.After(a.config.ShutdownTimeout / 2)

	for {
		select {
		case auditLog := <-a.logChan:
			if auditLog == nil {
				flush()
				return
			}
			*batch = append(*batch, auditLog)
			if len(*batch) >= a.config.BatchSize {
				flush()
			}
		case <-drainTimeout:
			a.log.Warn("Worker drain timed out",
				zap.Int("worker_id", workerID),
				zap.Int("remaining_in_batch", len(*batch)),
				zap.Int("remaining_in_queue", len(a.logChan)),
			)
			flush()
			return
		default:
			if len(a.logChan) == 0 {
				flush()
				return
			}
		}
	}
}
