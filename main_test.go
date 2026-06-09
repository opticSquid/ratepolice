package ratepolice_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/opticSquid/ratepolice"
	"github.com/opticSquid/ratepolice/shared"
)

var inMemoryFixedWindowConfig = shared.Config{
	MaxAllowedRequests: 10,
	KeyFunc: func(r *http.Request) string {
		return r.Header.Get("client-id")
	},
	TimeWindow:         1 * time.Minute,
	CoolDownDur:        2 * time.Second,
	CoolDownMultiplier: 3,
	Backend:            shared.InMemory,
	Algorithm:          shared.FixedWindowCounter,
	RedisConn:          "",
}

// What is the baseline cost per request for a single client?
func BenchmarkAllow_InMemory_FixedWindowCounter_SingleClient(b *testing.B) {
	rl, err := ratepolice.NewRateLimiter(inMemoryFixedWindowConfig)
	if err != nil {
		b.Fatal(err)
	}
	clientId := uuid.New().String()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("client-id", clientId)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := rl.Limit(next)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
	b.StopTimer()
	rl.Stop()
}

// How does performance change as the number of independent clients grows?
func BenchmarkAllow_InMemory_FixedWindowCounter_ClientScaling(b *testing.B) {
	clientCounts := []int{10, 50, 100, 500}
	reqsPerClient := []int{3, 6, 10}

	for _, numClients := range clientCounts {
		for _, reqsPC := range reqsPerClient {
			b.Run(fmt.Sprintf("clients=%d_reqsPerClient=%d", numClients, reqsPC), func(b *testing.B) {
				rl, err := ratepolice.NewRateLimiter(inMemoryFixedWindowConfig)
				if err != nil {
					b.Fatal(err)
				}
				clientIds := make([]string, numClients)
				for i := range clientIds {
					clientIds[i] = uuid.New().String()
				}
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
				handler := rl.Limit(next)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					clientIdx := (i / reqsPC) % numClients
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set("client-id", clientIds[clientIdx])
					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)
				}
				b.StopTimer()
				rl.Stop()
			})
		}
	}
}

// How does performance change as concurrent requests increase?
func BenchmarkAllow_InMemory_FixedWindowCounter_ConcurrencyScaling(b *testing.B) {
	parallelism := []int{1, 4, 16, 64, 256}
	numClients := 100

	for _, p := range parallelism {
		b.Run(fmt.Sprintf("parallel=%d", p), func(b *testing.B) {
			rl, err := ratepolice.NewRateLimiter(inMemoryFixedWindowConfig)
			if err != nil {
				b.Fatal(err)
			}
			clientIds := make([]string, numClients)
			for i := range clientIds {
				clientIds[i] = uuid.New().String()
			}
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			handler := rl.Limit(next)

			b.SetParallelism(p)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					clientIdx := i % numClients
					i++
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set("client-id", clientIds[clientIdx])
					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)
				}
			})
			b.StopTimer()
			rl.Stop()
		})
	}
}

// How does the full combination of clients + concurrency interact?
func BenchmarkAllow_InMemory_FixedWindowCounter_LoadMatrix(b *testing.B) {
	cases := []struct {
		clients       int
		reqsPerClient int
		parallelism   int
	}{
		{10, 5, 1},
		{10, 10, 16},
		{100, 5, 16},
		{100, 10, 64},
		{500, 5, 64},
		{500, 10, 256},
	}

	for _, tc := range cases {
		b.Run(fmt.Sprintf("clients=%d_rpc=%d_parallel=%d", tc.clients, tc.reqsPerClient, tc.parallelism), func(b *testing.B) {
			rl, err := ratepolice.NewRateLimiter(inMemoryFixedWindowConfig)
			if err != nil {
				b.Fatal(err)
			}
			clientIds := make([]string, tc.clients)
			for i := range clientIds {
				clientIds[i] = uuid.New().String()
			}
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			handler := rl.Limit(next)

			b.SetParallelism(tc.parallelism)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					clientIdx := (i / tc.reqsPerClient) % tc.clients
					i++
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set("client-id", clientIds[clientIdx])
					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)
				}
			})
			b.StopTimer()
			rl.Stop()
		})
	}
}
