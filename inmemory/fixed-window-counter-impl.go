package inmemory

import (
	"errors"
	"net/http"
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

func (cfg *InMemoryRateLimiter) fixedWindowCounter(clientId string) (shared.ResponseHeaders, error) {
	if rqCount, ok := cfg.data[clientId]; ok {
		if rqCount > cfg.maxAllowedRequests {
			return shared.ResponseHeaders{
				Type:                 "Rate Limit Exceeded",
				Title:                http.StatusText(http.StatusTooManyRequests),
				Detail:               "Rate limit has been exceeded",
				XRatelimitLimit:      shared.Blocked,
				XRatelimitRemaining:  0,
				XRatelimitReset:      cfg.windowEnd,
				XRatelimitRetryAfter: time.Until(cfg.windowEnd),
			}, errors.New("Request Blocked due to rate limit being crossed")
		}
		return shared.ResponseHeaders{
			XRatelimitLimit:     shared.Allowed,
			XRatelimitRemaining: cfg.maxAllowedRequests - rqCount,
			XRatelimitReset:     cfg.windowEnd,
		}, nil
	}
	cfg.data[clientId] = 1
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
