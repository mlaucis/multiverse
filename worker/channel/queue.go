/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package channel defines the logic for a queue and worker that are based on go channels
package channel

import (
	"sync/atomic"

	"github.com/tapglue/backend/worker"
)

type (
	queue struct {
		channel chan string
		size    int32
	}
)

func (q *queue) Add(message string) error {
	q.channel <- message
	atomic.AddInt32(&q.size, 1)
	return nil
}

func (q *queue) Get() (msg string, err error) {
	msg = <-q.channel
	atomic.AddInt32(&q.size, -1)
	return
}

func (q *queue) Size() int32 {
	return q.size
}

// NewQueue returns a new go channels queue
func NewQueue() worker.Queue {
	return &queue{
		channel: make(chan string),
	}
}
