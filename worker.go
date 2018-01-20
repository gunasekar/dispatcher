package dispatcher

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// Worker represents the worker that executes the job
type Worker struct {
	WorkerID        string
	WorkerPool      chan chan Job
	JobChannel      chan Job
	shutdown        chan bool
	confirmShutdown chan bool
}

// NewWorker ...
func NewWorker(workerID string, workerPool chan chan Job) *Worker {
	if workerPool == nil {
		log.WithFields(module).Errorf("WorkerPool channel is nil/not initialized")
		return nil
	}

	if workerID == "" {
		log.WithFields(module).Errorf("WorkerID is not set")
		return nil
	}

	return &Worker{
		WorkerID:        workerID,
		WorkerPool:      workerPool,
		JobChannel:      make(chan Job),
		shutdown:        make(chan bool, 1),
		confirmShutdown: make(chan bool, 1),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				log.WithFields(module).Debugf("%v: Received a job", w.WorkerID)
				w.executeJob(job)
			case <-w.shutdown:
				// we have received a signal to stop
				log.WithFields(module).Debugf("%v: Quitting the worker", w.WorkerID)
				w.confirmShutdown <- true
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	// Sending the signal to shutdown the worker
	w.shutdown <- true

	// Wait for the confirmation
	<-w.confirmShutdown
}

func (w *Worker) executeJob(job Job) {
	defer func() {
		if e := recover(); e != nil {
			log.WithFields(module).Errorf("%v: Job execution failure: %s", w.WorkerID, debug.Stack())
		}
	}()

	if errors := job.DoJob(); errors != nil && (len(errors) > 0) {
		log.WithFields(module).Debugf(w.WorkerID, "Worker %v: Error/s in doing the job. Job Object: %v", w.WorkerID, job)
		for i := range errors {
			if errors[i] != nil {
				log.WithFields(module).Errorf("%v: Error  %v: %v", w.WorkerID, i, errors[i].Error())
			}
		}
	}
}
