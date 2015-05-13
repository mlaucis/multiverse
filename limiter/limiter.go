/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package limiter takes care of limiting HTTP requests (and not only)
package limiter

import "time"

type (
	Limitee struct {
		// Hash returns a consistent hash for which the limit applies
		Hash string

		// Limit holds the number of requests one is allowed to do
		Limit int64

		// WindowSize holds the number of seconds for which the rate limit applies
		WindowSize int64
	}

	// Limiter is the actual the one providing the actual limitation implementation
	Limiter interface {

		// Request accepts a limitee parameter and for that it checks if it's still within
		// the limits or not. If not, it will return -1. If yes, it will decrement the remaining
		// number of hits by 1.
		Request(*Limitee) (int64, time.Time, error)
	}

	limitee struct {
		hash       string
		limit      int64
		windowSize int64
	}
)
