package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"golang.org/x/time/rate"
)

var ErrTooManyRequests = errors.New("number of permitted requests exceeded")

var rateLimit *RateLimiter

func RateLimit() *RateLimiter {
	if rateLimit == nil {
		OneThousandReqPS()
	}
	return rateLimit
}

type RateLimiter struct {
	l *rate.Limiter
}

var _ ratelimit.Limiter = (*RateLimiter)(nil)

func (r *RateLimiter) Limiter() *rate.Limiter {
	return r.l
}

func (r *RateLimiter) Limit(_ context.Context) error {
	if r.l.Allow() {
		return nil
	}
	return ErrTooManyRequests
}

func NReqPS(n int) *RateLimiter {
	return NewRateLimiter(time.Millisecond, n)
}
func OneThousandReqPS() *RateLimiter {
	return NReqPS(1000)
}
func TwoThousandReqPS() *RateLimiter {
	return NReqPS(2000)
}

func NewRateLimiter(s time.Duration, t int) *RateLimiter {
	rateLimit = &RateLimiter{
		l: rate.NewLimiter(rate.Every(s), t),
	}
	return rateLimit
}
