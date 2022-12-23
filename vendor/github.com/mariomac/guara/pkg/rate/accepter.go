package rate

import (
	"math"
	"time"
)

var clock = time.Now

// Accepter can accept or reject the arrival of a given event (e.g. a service
// request) based on a limit of requests per time duration. It uses a simple
// linear function to limit the rate of events so it's very efficient in terms
// of CPU and memory (doesn't need to store buckets nor using sliding windows
// for calculations)
type Accepter struct {
	// slope (negative) of the rate limiting rect
	slope    float64
	maxReqs  float64
	requests float64
	lastReq  time.Time
}

// NewAccepter returns an accepter that is able to accept a maximum number of
// maxRequests each perTime period.
func NewAccepter(maxRequests float64, perTime time.Duration) *Accepter {
	return &Accepter{
		maxReqs: maxRequests,
		slope:   -maxRequests / float64(perTime),
		lastReq: time.Now(),
	}
}

// Accept returns whether an incoming event should be accepted (it is within the
// given rate limits) or denied (frequency of events has been lately over the maximum
// allowed). If the event is accepted, this method internally updates the status of
// the instance for future limitation).
// This method is not thread-safe.
func (l *Accepter) Accept() bool {
	now := clock()
	timeDelta := now.Sub(l.lastReq)
	update := math.Max(l.requests+l.slope*float64(timeDelta), 0) + 1
	if update > l.maxReqs {
		return false
	}
	l.lastReq = now
	l.requests = update
	return true
}
