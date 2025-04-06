package main

import (
	"fmt"
	"net/http"
	ratelimiter "thing/rate-limiter/rate_limiter"
	"time"

	"github.com/gorilla/mux"
)

// test the rate limiter
func main() {
	RateLimit := ratelimiter.MakeRateLimiter(time.Minute*2, 5)
	_ = RateLimit

	r := mux.NewRouter()
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("key")

		if key == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}
		valid := RateLimit(r.Header.Get("key"))

		if !valid {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	srv := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}

	fmt.Println("listening on port :8080")
	srv.ListenAndServe()
}
