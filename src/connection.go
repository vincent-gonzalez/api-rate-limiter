package main

import (
	"time"
	"golang.org/x/time/rate"
)

// Connection defines a request made the API
type Connection struct {
	Identifier string
	Limiter *rate.Limiter
	LastActive time.Time
}
