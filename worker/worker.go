/**
 * @author Florin Patan <florinpatan@gmail.com
 */

// Package worker holds the worker queues and the workers
package worker

type (
	// Status defines the status of a worker
	Status int

	// Queue defines the methods a message queue must have
	Queue interface {
		Add(string) error
		Get() (string, error)
		Size() int32
	}

	// Worker defines the methods a worker must have
	Worker interface {
		Process()
		Start() error
		Stop() error
		Status() Status
	}
)

const (
	// Stopped means that the worker is stopped
	Stopped = Status(iota)

	// Started means that the worker is running
	Started
)
