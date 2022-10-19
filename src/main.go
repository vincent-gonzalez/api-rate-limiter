package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := struct{
		CurrentTime string `json:"currentTime"`
	}{
		CurrentTime: time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(currentTime)
}

func rateLimitMiddleware(next http.Handler, rateLimiter *RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !rateLimiter.IsRequestAllowed(r.RemoteAddr) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Starting server...")

	rateLimitInterval := flag.Int("rateLimit", 1, "rate limit in milliseconds")
	requestLimit := flag.Int("requestLimit", 1, "max number of requests per rate interval")
	flag.Parse()

	rateLimiter := NewRateLimiter(*rateLimitInterval, *requestLimit)
	stop := rateLimiter.RunInactiveConnectionsCleanupRoutine(1, time.Minute, 2, time.Minute)
	defer close(stop)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", rateLimitMiddleware(mux, rateLimiter))
	if err != nil {
		log.Fatalf("server failed: %v", err.Error())
	}
}
