package inmemory

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

func NewInMemoryRateLimiter(cfg shared.Config) *InMemoryRateLimiter {
	ctx, cncl := context.WithCancel(context.Background())
	curTime := time.Now()
	rl := &InMemoryRateLimiter{
		algorithm:          cfg.Algorithm,
		maxAllowedRequests: cfg.MaxAllowedRequests,
		keyFunc:            cfg.KeyFunc,
		timeWindow:         cfg.TimeWindow,
		coolDownDur:        cfg.CoolDownDur,
		coolDownMultiplier: cfg.CoolDownMultiplier,
		windowStart:        curTime,
		windowEnd:          curTime.Add(cfg.TimeWindow),
		ctx:                ctx,
		cnclFunc:           cncl,
		data:               make(map[string]int64),
	}
	if rl.timeWindow > 0 {
		switch rl.algorithm {
		case shared.FixedWindowCounter:
			go rl.fixedWindowCleanup()
		}
	}
	return rl
}

func (cfg *InMemoryRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Ratelimiting process in memory
		response, err := cfg.algoSwitcher(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		// handing over to next middleware
		next.ServeHTTP(w, r)
	})
}

func (cfg *InMemoryRateLimiter) Stop() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.cnclFunc()
}

func (cfg *InMemoryRateLimiter) algoSwitcher(w http.ResponseWriter, r *http.Request) (shared.ResponseHeaders, error) {
	clientId := cfg.keyFunc(r)
	switch cfg.algorithm {
	case shared.FixedWindowCounter:
		return cfg.fixedWindowCounter(clientId)
	case shared.SlidingWindow:
	case shared.SlidingWindowCounter:
	case shared.TokenBucket:
	case shared.LeakyBucket:
	}
	//dummy will never come to this
	return shared.ResponseHeaders{}, nil
}
