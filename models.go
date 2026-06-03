package ratepolice

import (
	"net/http"
	"time"
)

// Rate Limiting Algorithm to be used
type Algorithm string

const (
	FixedWindowCounter   Algorithm = "fixed_window_counter"
	SlidingWindow        Algorithm = "sliding_window"
	SlidingWindowCounter Algorithm = "sliding_window_counter"
	TokenBucket          Algorithm = "token_bucket"
	LeakyBucket          Algorithm = "leaky_bucket"
)

// Backend to be used to store data
type Backend string

const (
	InMemory Backend = "in_memory"
	Redis    Backend = "redis"
)

type Config struct {
	Algoritm        Algorithm
	Backend         Backend
	AllowedRequests int64
	KeyFunc         func(args ...any) string
	TimeWindow      time.Duration
}

type RateLimiter struct {
	data            map[string][]time.Time
	algoritm        Algorithm
	backend         Backend
	allowedRequests int64
	keyFunc         func(args ...any) string
	timeWindow      time.Duration
}

type Verdict string

const (
	Allowed Verdict = "allowed"
	Blocked Verdict = "blocked"
)

type ResponseHeaders struct {
	XRatelimitLimit     Verdict       `json:"X-RateLimit-Limit"`
	XRatelimitRemaining int64         `json:"X-RateLimit-Remaining"`
	XRatelimitReset     time.Time     `json:"X-RateLimit-Reset"`
	RetryAfter          time.Duration `json:"Retry-After"`
}

// RFC 9457 compliant ProblemDetails Object
type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

type RatePolice interface {
	Limit(next http.Handler) http.Handler
}
