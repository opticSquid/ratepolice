package inmemory

import (
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
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
	for i := range 5 {
		slog.Info("iteration", "rq no", i)

		res, _ := inmemoryRl.fixedWindowCounter(client_id)

		if res.XRatelimitLimit != shared.Allowed {
			t.Errorf("expected allowed, got %s, rq no: %d", res.XRatelimitLimit, i)
		}
		if res.XRatelimitRemaining != 5-i {
			t.Errorf("expected remaining %d, got %d, rq no: %d", 5-i, res.XRatelimitRemaining, i)
		}
	}
}

func Test_InMemoryRateLimiter_Single_Client_Exceeding_Limit(t *testing.T) {
	client_id := "test-client-0"

	inmemoryRl := NewInMemoryRateLimiter(config)
	slog.Info("doing 4 requests to consume the rate limit")
	for i := range 5 {
		slog.Info("Sending Request", "rq no", i)
		res, _ := inmemoryRl.fixedWindowCounter(client_id)

		if res.XRatelimitLimit != shared.Allowed {
			t.Errorf("expected allowed, got %s, rq no: %d", res.XRatelimitLimit, i)
		}
		if res.XRatelimitRemaining != 5-i {
			t.Errorf("expected remaining %d, got %d, rq no: %d", 5-i, res.XRatelimitRemaining, i)
		}
	}
	res, _ := inmemoryRl.fixedWindowCounter(client_id)
	if res.XRatelimitLimit != shared.Blocked {
		t.Errorf("expected blocked, got %s", res.XRatelimitLimit)
	}
	if res.XRatelimitRemaining != 0 {
		t.Errorf("expected remaining 0, got %d", res.XRatelimitRemaining)
	}
}

func randomJitter() {
	// 5–30ms random delay
	time.Sleep(time.Duration(5+rand.Intn(25)) * time.Millisecond)
}

func Test_InMemoryRateLimiter_Multiple_Clients_Within_Limit_Mixed_Sequence(t *testing.T) {
	inmemoryRl := NewInMemoryRateLimiter(config)

	clients := []string{
		"client-1",
		"client-2",
		"client-3",
		"client-4",
		"client-5",
	}

	var wg sync.WaitGroup

	for _, clientID := range clients {
		wg.Add(1)

		// Parallel across clients
		go func(cid string) {
			defer wg.Done()

			// Sequential per client
			for i := 0; i < int(config.MaxAllowedRequests); i++ {
				randomJitter()
				slog.Info("Sending Request", "client", cid, "rq no", i)

				res, _ := inmemoryRl.fixedWindowCounter(cid)

				if res.XRatelimitLimit != shared.Allowed {
					t.Errorf("expected allowed for client %s, got %s at rq %d",
						cid, res.XRatelimitLimit, i)
				}

				expectedRemaining := config.MaxAllowedRequests - i
				if res.XRatelimitRemaining != expectedRemaining {
					t.Errorf("client %s: expected remaining %d, got %d",
						cid, expectedRemaining, res.XRatelimitRemaining)
				}
			}
		}(clientID)
	}

	wg.Wait()
}

func Test_InMemoryRateLimiter_Multiple_Clients_Exceeding_Limit_Mixed_Sequence(t *testing.T) {
	inmemoryRl := NewInMemoryRateLimiter(config)

	clients := []string{
		"client-1",
		"client-2",
		"client-3",
		"client-4",
		"client-5",
		"client-6",
	}

	var wg sync.WaitGroup

	for _, clientID := range clients {
		wg.Add(1)

		go func(cid string) {
			defer wg.Done()

			slog.Info("Testing client", "client", cid)

			// Sequential: consume quota
			for i := 0; i < int(config.MaxAllowedRequests); i++ {
				res, _ := inmemoryRl.fixedWindowCounter(cid)

				if res.XRatelimitLimit != shared.Allowed {
					t.Errorf("expected allowed for client %s, got %s at rq %d",
						cid, res.XRatelimitLimit, i)
				}
			}
			randomJitter()
			// Extra request — should be blocked
			res, _ := inmemoryRl.fixedWindowCounter(cid)

			if res.XRatelimitLimit != shared.Blocked {
				t.Errorf("expected blocked for client %s, got %s",
					cid, res.XRatelimitLimit)
			}

			if res.XRatelimitRemaining != 0 {
				t.Errorf("expected remaining 0 for client %s, got %d",
					cid, res.XRatelimitRemaining)
			}
		}(clientID)
	}

	wg.Wait()
}
