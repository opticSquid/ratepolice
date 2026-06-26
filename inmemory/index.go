package inmemory

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	}
	switch rl.algorithm {
	case shared.FixedWindowCounter:
		rl.fixedWindowCounterData = make(map[string]int)
	case shared.SlidingWindowCounter:
		rl.slidingWindowData = make(map[string][]time.Time)
	}
	if rl.timeWindow > 0 {
		switch rl.algorithm {
		case shared.FixedWindowCounter:
			go rl.fixedWindowCleanup()
		case shared.SlidingWindowCounter:
			go rl.slidingWindowCleanup()
		}
	}
	return rl
}

func setRateLimitHeaders(w http.ResponseWriter, h shared.ResponseHeaders) {
	w.Header().Set("X-Ratelimit-Limit", string(h.XRatelimitLimit))
	w.Header().Set("X-Ratelimit-Remaining", strconv.Itoa(h.XRatelimitRemaining))
	w.Header().Set("X-Ratelimit-Reset", h.XRatelimitReset.UTC().Format(time.RFC3339))
	if h.XRatelimitRetryAfter > 0 {
		w.Header().Set("X-Ratelimit-Retry-After", strconv.FormatInt(int64(h.XRatelimitRetryAfter), 10))
	}
}

func (rl *InMemoryRateLimiter) getBlockedResponseHeaders() shared.ResponseHeaders {
	return shared.ResponseHeaders{
		XRatelimitLimit:      shared.Blocked,
		XRatelimitRemaining:  0,
		XRatelimitReset:      rl.windowEnd,
		XRatelimitRetryAfter: rl.windowEnd.Sub(time.Now()),
	}
}
func getBlockedResponseBody(title, reason string) shared.ErrorResponseBody {
	return shared.ErrorResponseBody{
		Type:   "Blocked",
		Title:  title,
		Detail: reason,
	}
}

func writeResponse(w http.ResponseWriter, status int, headers shared.ResponseHeaders, body shared.ErrorResponseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("ratepolice: failed to encode response body: %v", err)
	}
}

func (rl *InMemoryRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.maxAllowedRequests == 0 {
			respH := rl.getBlockedResponseHeaders()
			respB := getBlockedResponseBody(http.StatusText(http.StatusTooManyRequests), "no request is allowed for this endpoint")
			setRateLimitHeaders(w, respH)
			writeResponse(w, http.StatusServiceUnavailable, respH, respB)
			return
		}

		respH, err := rl.algoSwitcher(w, r)
		if err != nil {
			respB := getBlockedResponseBody(http.StatusText(http.StatusTooManyRequests), "too many requests")
			setRateLimitHeaders(w, respH)
			writeResponse(w, http.StatusTooManyRequests, respH, respB)
			return
		}

		setRateLimitHeaders(w, respH)
		next.ServeHTTP(w, r)
	})
}

func (rl *InMemoryRateLimiter) Stop() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.cnclFunc()
}

func (rl *InMemoryRateLimiter) algoSwitcher(w http.ResponseWriter, r *http.Request) (shared.ResponseHeaders, error) {
	clientId := rl.keyFunc(r)
	switch rl.algorithm {
	case shared.FixedWindowCounter:
		return rl.fixedWindowCounter(clientId)
	case shared.SlidingWindow:
		return rl.slidingWindow(clientId)
	case shared.SlidingWindowCounter:
	case shared.TokenBucket:
	case shared.LeakyBucket:
	}
	//dummy will never come to this
	return shared.ResponseHeaders{}, nil
}
