package inmemory

import (
	"net/http"
	"testing"
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

var config = shared.Config{
	MaxAllowedRequests: 5,
	KeyFunc: func(r *http.Request) string {
		return r.Header.Get("client-id")
	},
	TimeWindow:         10 * time.Second,
	CoolDownDur:        2 * time.Second,
	CoolDownMultiplier: 3,
	Backend:            shared.InMemory,
	Algorithm:          shared.FixedWindowCounter,
	RedisConn:          "",
}

func Test_InMemoryRateLimiter_Single_Client_Within_Limit(t *testing.T) {
	client_id := "test-client-0"
	inmemoryRl := NewInMemoryRateLimiter(config)
	for i := 0; i < 5; i++ {
		res, _ := inmemoryRl.fixedWindowCounter(client_id)
		if res.XRatelimitLimit != shared.Allowed {
			t.Errorf("expected allowed, got %s", res.XRatelimitLimit)
		}
		if res.XRatelimitRemaining != int64(5-i) {
			t.Errorf("expected remaining %d, got %d", 5-i, res.XRatelimitRemaining)
		}
	}
}
