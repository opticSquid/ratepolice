package inmemory

import (
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

func (cfg *InMemoryRateLimiter) getOrCreate(clientId string) int64 {
	cfg.mu.Lock() // Allows Concurrent Reads
	defer cfg.mu.Unlock()
	if rqCount, ok := cfg.data[clientId]; ok {
		return rqCount
	}
	cfg.data[clientId] = 0
	return 0
}

func (cfg *InMemoryRateLimiter) fixedWindowCounter(clientId string) (shared.ResponseHeaders, error) {
	rqCount := cfg.getOrCreate(clientId)
	if cfg.maxAllowedRequests != 0 && rqCount > cfg.maxAllowedRequests {
		return shared.ResponseHeaders{
			XRatelimitLimit:     shared.Allowed,
			XRatelimitRemaining: cfg.maxAllowedRequests - rqCount,
			XRatelimitReset:     cfg.windowEnd,
		}, nil
	}
	cfg.mu.Lock()
	cfg.data[clientId] = rqCount + 1
	cfg.mu.Unlock()
	return shared.ResponseHeaders{
		XRatelimitLimit:     shared.Allowed,
		XRatelimitRemaining: cfg.maxAllowedRequests - 1,
		XRatelimitReset:     cfg.windowEnd,
	}, nil
}
func (cfg *InMemoryRateLimiter) fixedWindowPurge() {
	curTime := time.Now()
	cfg.windowStart = curTime
	cfg.windowEnd = curTime.Add(cfg.timeWindow)
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	for i := range cfg.data {
		cfg.data[i] = 0
	}
}
func (cfg *InMemoryRateLimiter) fixedWindowCleanup() {
	ticker := time.NewTicker(cfg.timeWindow)
	for {
		select {
		case <-ticker.C:
			cfg.fixedWindowPurge()
		case <-cfg.ctx.Done():
			return
		}
	}
}
