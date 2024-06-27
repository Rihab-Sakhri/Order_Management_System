package dispatcher

import (
	"fmt"

	"github.com/sofiukl/oms/oms-checkout/core"
	"github.com/sofiukl/oms/oms-checkout/worker"
)

// StartDispatcher - dispatcher start
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue := make(chan chan worker.Work, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := worker.NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-core.WorkQueue:
				fmt.Println("Received work requeust", work)
				go func() {
					// getting available work channel
					workChannel := <-WorkerQueue

					// dispatching the work
					fmt.Println("Dispatching work request")
					workChannel <- work
				}()
			}
		}
	}()
}
