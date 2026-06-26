package inmemory

import (
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

func (rl *InMemoryRateLimiter) getOrCreateClientFixedWindowCounter(clientId string) int {
	rl.mu.Lock() // Allows Concurrent Reads
	defer rl.mu.Unlock()
	if rqCount, ok := rl.fixedWindowCounterData[clientId]; ok {
		return rqCount
	}
	rl.fixedWindowCounterData[clientId] = 0
	return 0
}

func (rl *InMemoryRateLimiter) fixedWindowCounter(clientId string) (shared.ResponseHeaders, error) {
	rqCount := rl.getOrCreateClientFixedWindowCounter(clientId)
	if rqCount >= rl.maxAllowedRequests {
		return rl.getBlockedResponseHeaders(), nil
	}

	remaining := rl.maxAllowedRequests - rqCount
	rl.mu.Lock()
	rl.fixedWindowCounterData[clientId] = rqCount + 1
	rl.mu.Unlock()
	return shared.ResponseHeaders{
		XRatelimitLimit:     shared.Allowed,
		XRatelimitRemaining: remaining,
		XRatelimitReset:     rl.windowEnd,
	}, nil
}
func (rl *InMemoryRateLimiter) fixedWindowPurge() {
	curTime := time.Now()
	rl.windowStart = curTime
	rl.windowEnd = curTime.Add(rl.timeWindow)
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for i := range rl.fixedWindowCounterData {
		rl.fixedWindowCounterData[i] = 0
	}
}
func (rl *InMemoryRateLimiter) fixedWindowCleanup() {
	ticker := time.NewTicker(rl.timeWindow)
	for {
		select {
		case <-ticker.C:
			rl.fixedWindowPurge()
		case <-rl.ctx.Done():
			return
		}
	}
}
