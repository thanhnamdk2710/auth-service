package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	DBConnectionsOpen *prometheus.GaugeFunc
	DBConnectionsInUse *prometheus.GaugeFunc
	DBConnectionsIdle *prometheus.GaugeFunc

	UserRegistrations prometheus.Counter
	LoginAttempts     *prometheus.CounterVec
	PasswordResets    prometheus.Counter
}

func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		HTTPRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestsInFlight: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),
		UserRegistrations: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Name: "user_registrations_total",
				Help: "Total number of user registrations",
			},
		),
		LoginAttempts: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "login_attempts_total",
				Help: "Total number of login attempts",
			},
			[]string{"status"},
		),
		PasswordResets: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Name: "password_resets_total",
				Help: "Total number of password reset requests",
			},
		),
	}

	return m
}

func (m *Metrics) RegisterDBStats(reg prometheus.Registerer, db *sql.DB) {
	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "db_connections_open",
			Help: "Number of open database connections",
		},
		func() float64 {
			return float64(db.Stats().OpenConnections)
		},
	))

	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
		func() float64 {
			return float64(db.Stats().InUse)
		},
	))

	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
		func() float64 {
			return float64(db.Stats().Idle)
		},
	))
}
