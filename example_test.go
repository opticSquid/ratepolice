package ratepolice_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/opticSquid/ratepolice"
)

func Example() {
	cfg := ratepolice.Config{
		MaxAllowedRequests: 100,
		KeyFunc:            func(req *http.Request) string { return req.RemoteAddr },
		TimeWindow:         1 * time.Minute,
		Backend:            ratepolice.InMemory,
	}
	rl, err := ratepolice.NewRateLimiter(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rl.Stop()

	handler := rl.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	_ = handler
	fmt.Println("rate limiter ready")
	// Output: rate limiter ready
}
