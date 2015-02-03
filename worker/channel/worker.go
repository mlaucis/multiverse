/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package channel

import "github.com/tapglue/backend/worker"

type (
	wrkr struct {
		queue  worker.Queue
		status worker.Status
	}
)

func (w *wrkr) Process() {
	if w.queue.Size() == 0 {
		return
	}

	message, err := w.queue.Get()
	if err == nil {
		// TODO oh well
		panic(err)
	}

	// TODO do something with the message here
	_ = message
}

func (w *wrkr) Start() error {
	w.status = worker.Started
	return nil
}

func (w *wrkr) Stop() error {
	w.status = worker.Stopped
	return nil
}

func (w *wrkr) Status() worker.Status {
	return w.status
}

// NewWorker creates a new worker for the specified queue
func NewWorker(queue worker.Queue) worker.Worker {
	return &wrkr{
		queue:  queue,
		status: worker.Stopped,
	}
}
