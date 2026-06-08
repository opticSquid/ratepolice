package inmemory

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

type InMemoryRateLimiter struct {
	algorithm          shared.Algorithm
	maxAllowedRequests int64
	keyFunc            func(*http.Request) string
	timeWindow         time.Duration
	coolDownDur        time.Duration
	coolDownMultiplier int
	mu                 sync.Mutex
	windowStart        time.Time
	windowEnd          time.Time
	ctx                context.Context
	cnclFunc           context.CancelFunc
	data               map[string]int64
}
