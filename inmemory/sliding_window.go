package inmemory

import (
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

func (rl *InMemoryRateLimiter) getOrCreateClientSlidingWindow(clientId string) []time.Time {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rqLog, ok := rl.slidingWindowData[clientId]; ok {
		return rqLog
	}
	rl.slidingWindowData[clientId] = make([]time.Time, 0, rl.maxAllowedRequests)
	return rl.slidingWindowData[clientId]
}

func (rl *InMemoryRateLimiter) slidingWindow(clientId string) (shared.ResponseHeaders, error) {
	// rqLog := rl.getOrCreateClientSlidingWindow(clientId)

	return shared.ResponseHeaders{}, nil
}

func (rl *InMemoryRateLimiter) slidingWindowCleanup() {

}
