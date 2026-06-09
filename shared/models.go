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
	//Maximum number of allowed requests per client with in the time window
	MaxAllowedRequests int64
	// The function to uniquely identify a client given a request
	KeyFunc func(*http.Request) string
	// The time window for rate limiting
	TimeWindow time.Duration
	// The duration of the cool down period after a blocked request
	CoolDownDur time.Duration
	// The multiplier for the cool down period after a blocked request
	CoolDownMultiplier int
	// The backend to be used to store data
	Backend Backend
	// The algorithm to be used for rate limiting
	Algorithm Algorithm
	// The Redis connection string to be used for storing data in case you are using Redis as the backend
	RedisConn string
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
