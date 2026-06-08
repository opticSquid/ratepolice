package shared

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
	MaxAllowedRequests int64
	KeyFunc            func(*http.Request) string
	TimeWindow         time.Duration
	CoolDownDur        time.Duration
	CoolDownMultiplier int
	Backend            Backend
	Algorithm          Algorithm
	RedisConn          string
}

type Verdict string

const (
	Allowed Verdict = "allowed"
	Blocked Verdict = "blocked"
)

// RFC 9457 ProblemDetails object extended ResponseHeaders Object
type ResponseHeaders struct {
	Type                 string        `json:"type,omitempty"`
	Title                string        `json:"title"`
	Detail               string        `json:"detail"`
	Instance             string        `json:"instance,omitempty"`
	XRatelimitLimit      Verdict       `json:"X-RateLimit-Limit"`
	XRatelimitRemaining  int64         `json:"X-RateLimit-Remaining"`
	XRatelimitReset      time.Time     `json:"X-RateLimit-Reset"`
	XRatelimitRetryAfter time.Duration `json:"X-RateLimit-Retry-After"`
}

type RatePolice interface {
	Limit(next http.Handler) http.Handler
	Stop()
}
