package main

import (
	"sync"
	"time"
	"golang.org/x/time/rate"
)

// RateLimiter limits the number of requests made to an HTTP server
type RateLimiter struct {
	IPConnections map[string]*Connection
	// RequestLimit is the maximum number of requests that can occur per time interval
	RequestLimit int
	// RateLimit is the time interval in which requests up to the
	// maximum limit can occur
	RateLimit time.Duration
	mutexLock sync.Mutex
}

// NewRateLimiter creates new instances of RateLimiter
func NewRateLimiter(interval int, count int) *RateLimiter {
	var rateLimit time.Duration

	if count > 0 {
		rateLimit = (time.Duration(interval)*time.Millisecond) / time.Duration(count)
	} else {
		rateLimit = 0
	}

	return &RateLimiter{
		IPConnections: make(map[string]*Connection),
		RequestLimit: count,
		RateLimit: rateLimit,
	}
}

// IsRequestAllowed determines if the rate limit for a connection has been reached
func (r *RateLimiter) IsRequestAllowed(clientIP string) bool {
	r.mutexLock.Lock()
	connection, hasConnection := r.IPConnections[clientIP]
	r.mutexLock.Unlock()

	if !hasConnection {
		r.TrackConnection(clientIP)
		r.mutexLock.Lock()
		connection = r.IPConnections[clientIP]
		r.mutexLock.Unlock()
	}

	if connection.Limiter.Allow() {
		return true
	}

	return false
}

func (r *RateLimiter) TrackConnection(connectionIP string) {
	r.mutexLock.Lock()
	defer r.mutexLock.Unlock()
	// don't track a connection again if it is already tracked
	if _, hasConnection := r.IPConnections[connectionIP]; !hasConnection {
		newLimiter := rate.NewLimiter(rate.Every(r.RateLimit), r.RequestLimit)
		newConnection := &Connection{
			Identifier: connectionIP,
			Limiter: newLimiter,
			LastActive: time.Now(),
		}
		r.IPConnections[connectionIP] = newConnection
	}
}

// RunInactiveConnectionsCleanupRoutine runs a go routine that periodically
// removes connections that have not been active for a period of time.
func (r *RateLimiter) RunInactiveConnectionsCleanupRoutine(
	routineInterval int,
	intervalUnit time.Duration,
	inactiveConnectionInterval int,
	inactiveConnectionIntervalUnit time.Duration) chan bool {
	stop := make(chan bool)
	ticker := time.NewTicker(time.Duration(routineInterval) * intervalUnit)
	cleanupInterval := time.Duration(inactiveConnectionInterval) * inactiveConnectionIntervalUnit

	go func() {
		for {
			select {
			case <- ticker.C:
				r.mutexLock.Lock()
				for k, v := range r.IPConnections {
					if time.Since(v.LastActive) > cleanupInterval {
						delete(r.IPConnections, k)
					}
				}
				r.mutexLock.Unlock()
			case <- stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}
