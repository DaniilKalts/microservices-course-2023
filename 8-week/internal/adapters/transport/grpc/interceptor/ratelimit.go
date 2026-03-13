package interceptor

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const authServicePrefix = "/auth.v1.AuthV1/"

// RateLimitConfig holds per-IP token-bucket parameters.
// Auth endpoints use a separate, tighter limit to protect against brute-force.
type RateLimitConfig struct {
	RPS   float64
	Burst int

	AuthRPS   float64
	AuthBurst int
}

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// perIPLimiter maintains a separate token bucket for each client IP.
type perIPLimiter struct {
	mu      sync.Mutex
	entries map[string]*ipEntry

	rps   rate.Limit
	burst int
}

func newPerIPLimiter(rps float64, burst int) *perIPLimiter {
	l := &perIPLimiter{
		entries: make(map[string]*ipEntry),
		rps:     rate.Limit(rps),
		burst:   burst,
	}

	go l.evictLoop()

	return l
}

func (l *perIPLimiter) allow(ip string) bool {
	now := time.Now()

	l.mu.Lock()
	e, ok := l.entries[ip]
	if !ok {
		e = &ipEntry{limiter: rate.NewLimiter(l.rps, l.burst)}
		l.entries[ip] = e
	}
	e.lastSeen = now
	l.mu.Unlock()

	return e.limiter.Allow()
}

// evictLoop removes entries idle for 10+ minutes to prevent unbounded map growth.
func (l *perIPLimiter) evictLoop() {
	const evictInterval = 1 * time.Minute
	const maxIdleTime = 10 * time.Minute

	ticker := time.NewTicker(evictInterval)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-maxIdleTime)

		l.mu.Lock()
		for ip, e := range l.entries {
			if e.lastSeen.Before(cutoff) {
				delete(l.entries, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimitInterceptor enforces per-IP rate limits with a stricter bucket for auth endpoints.
func RateLimitInterceptor(cfg RateLimitConfig) grpc.UnaryServerInterceptor {
	defaultLimiter := newPerIPLimiter(cfg.RPS, cfg.Burst)
	authLimiter := newPerIPLimiter(cfg.AuthRPS, cfg.AuthBurst)

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		// Strip port — peer.Addr returns "ip:port".
		ip, _, err := net.SplitHostPort(p.Addr.String())
		if err != nil {
			return handler(ctx, req)
		}

		limiter := defaultLimiter
		if strings.HasPrefix(info.FullMethod, authServicePrefix) {
			limiter = authLimiter
		}

		if !limiter.allow(ip) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}
