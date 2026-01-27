package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	CleanupInterval   time.Duration
	TTL               time.Duration
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute,
		TTL:               time.Minute * 5,
	}
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*client
	config  RateLimiterConfig
	stopCh  chan struct{}
}

func newRateLimiter(cfg RateLimiterConfig) *rateLimiter {
	rl := &rateLimiter{
		clients: make(map[string]*client),
		config:  cfg,
		stopCh:  make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	c, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if exists {
		rl.mu.Lock()
		c.lastSeen = time.Now()
		rl.mu.Unlock()
		return c.limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if c, exists := rl.clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
	rl.clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, c := range rl.clients {
				if time.Since(c.lastSeen) > rl.config.TTL {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *rateLimiter) stop() {
	close(rl.stopCh)
}

func RateLimit(cfg RateLimiterConfig) gin.HandlerFunc {
	limiter := newRateLimiter(cfg)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.getLimiter(ip)

		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

func RateLimitDefault() gin.HandlerFunc {
	return RateLimit(DefaultRateLimiterConfig())
}
