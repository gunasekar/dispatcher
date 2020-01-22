package dispatcher

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// LocalDispatcher - Go channel based job dispatcher for managed job handling
// where the jobs will be produced and consumed within the instance
type LocalDispatcher struct {
	Name            string
	WorkerPool      chan chan Job
	MaxWorkers      int
	JobQueue        chan Job
	EnableDebugLogs bool
	workers         []*Worker
	shutdown        chan bool
	confirmShutdown chan bool
}

// NewLocalDispatcher - Dispaches the configures number of workers
// based on the jobs available from the provided Go Channel (JobQueue)
// Jobs should be produced to the assigned go channel (JobQueue)
func NewLocalDispatcher(name string,
	maxWorkers int,
	jobQueue chan Job,
) *LocalDispatcher {

	if name == "" {
		id := newUUID()
		name = "LocalDispatcher-" + id[0:5]
		log.WithFields(module).Infof(
			"Dispatcher name provided is nil/empty. Hence assigning a system generated name - %v", name)
	}

	if maxWorkers <= 0 {
		log.WithFields(module).Errorf("%v: Number of workers should be at least 1 for requested dispatcher. Given MaxWorkers: %v",
			name, maxWorkers)
		return nil
	}

	if jobQueue == nil {
		log.WithFields(module).Errorf("%v: Provided Job Queue is null for dispatcher", name)
		return nil
	}

	log.WithFields(module).Infof(
		"Creating Local Dispatcher - %v. MaxWorkers: %v",
		name, maxWorkers)

	workerPool := make(chan chan Job, maxWorkers)

	return &LocalDispatcher{
		Name:            name,
		WorkerPool:      workerPool,
		JobQueue:        jobQueue,
		MaxWorkers:      maxWorkers,
		shutdown:        make(chan bool, 1),
		confirmShutdown: make(chan bool, 1),
	}
}

// GetName - Gets the name of the dispatcher
func (d *LocalDispatcher) GetName() string {
	return d.Name
}

// Run - Dispatches the workers as per the configuration
func (d *LocalDispatcher) Run() {
	run(d, d.MaxWorkers, d.WorkerPool)
}

// Shutdown - Gracefully shuts down the dispatcher
// by closing the job enqueuing channel
// and waits for workers to complete the accepted jobs
func (d *LocalDispatcher) Shutdown() <-chan struct{} {
	// Closing the job queue to stop receiving the requests
	close(d.JobQueue)
	log.WithFields(module).Infof(
		"%v: Closed the job channel of the Local Dispatcher",
		d.Name)

	log.WithFields(module).Infof(
		"%v: Shutting down the Local Dispatcher",
		d.Name)

	// Sending the signal to shutdown the dispatcher
	d.shutdown <- true

	// Initializing the channel for the shutdown completion
	shutdownComplete := make(chan struct{})

	// Wait for all the pending jobs in the queue to get served
	<-d.confirmShutdown

	wg := new(sync.WaitGroup)
	wg.Add(len(d.workers))

	for i := range d.workers {
		go func(w *Worker, wg *sync.WaitGroup) {
			defer wg.Done()
			w.Stop()
		}(d.workers[i], wg)
	}

	wg.Wait()
	close(shutdownComplete)

	log.WithFields(module).Infof(
		"%v: Gracefully shut down the Local Dispatcher",
		d.Name)
	return shutdownComplete
}

func (d *LocalDispatcher) addWorker(w *Worker) {
	d.workers = append(d.workers, w)
}

func (d *LocalDispatcher) dispatch() {
	for {
		select {
		case job, ok := <-d.JobQueue:
			if ok == false {
				log.WithFields(module).Infof("%v: Job Queue Channel of the dispatcher is closed", d.Name)
				continue
			}

			// Job request has been received
			log.WithFields(module).Debugf("%v: Received a job with ID %v", d.Name, job.GetJobID())

			// Try to obtain a worker job channel that is available.
			// This will block until a worker is idle
			jobChannel := <-d.WorkerPool

			// Dispatch the job to the corresponding worker's job channel
			jobChannel <- job
		case <-d.shutdown:
			// We have received a signal to stop
			// No jobs in the queue. Confirm the dispatcher shutdown
			d.confirmShutdown <- true
			return
		}
	}
}
