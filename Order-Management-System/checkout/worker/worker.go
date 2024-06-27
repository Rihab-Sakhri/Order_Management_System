package worker

import (
	"log"

	"github.com/sofiukl/oms/oms-checkout/api"
)

// Worker - This is Worker struct
type Worker struct {
	ID          int
	WorkChannel chan Work
	WorkerQueue chan chan Work
}

// NewWorker - This creates the instance of new worker
func NewWorker(id int, workerQueue chan chan Work) Worker {
	worker := Worker{
		ID:          id,
		WorkChannel: make(chan Work),
		WorkerQueue: workerQueue,
	}

	return worker
}

// Start - This is the runnable method of the worker
func (w *Worker) Start() {
	go func() {
		for {
			// assigning available channel to WorkerQueue
			w.WorkerQueue <- w.WorkChannel
			select {
			case job := <-w.WorkChannel:
				// Receive a work request
				log.Printf("worker %d: working on %s!\n", w.ID, job.Work.CartID)
				api.CheckoutProduct(job.Conn, job.Config, job.Work, job.Lock)
				log.Printf("worker%d: work %s completed\n", w.ID, job.Work.CartID)
			}
		}
	}()
}
