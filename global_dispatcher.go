package dispatcher

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// GlobalDispatcher - Job Dispatcher for managed job handling in a distributed
// environment where the jobs will be produced and consumed globally
// JobConsumer object is needed which defines how to consume the job from the
// global queue/pipe in the distributed system (Check the JobConsumer interface)
type GlobalDispatcher struct {
	Name                     string
	WorkerPool               chan chan Job
	MaxWorkers               int
	JobConsumer              JobConsumer
	PollingIntervalInSeconds int
	workers                  []*Worker
	shutdown                 chan bool
	confirmShutdown          chan bool
}

// NewGlobalDispatcher - Dispaches the configures number of workers
// based on the jobs available from Global Queue in distributed environment
func NewGlobalDispatcher(name string,
	maxWorkers int,
	jobConsumer JobConsumer,
	pollingIntervalInSeconds int) *GlobalDispatcher {

	if name == "" {
		id := newUUID()
		name = "GlobalDispatcher-" + id[0:5]
		log.WithFields(module).Infof(
			"Dispatcher name provided is nil/empty. Hence assigning a system generated name - %v", name)
	}

	if maxWorkers <= 0 {
		log.WithFields(module).Infof("%v: Number of workers should be at least 1 for requested dispatcher. Given MaxWorkers: %v",
			name, maxWorkers)
		return nil
	}

	if jobConsumer == nil {
		log.WithFields(module).Errorf("%v: JobConsumer object is nil for dispatcher", name)
		return nil
	}

	if pollingIntervalInSeconds <= 0 {
		log.WithFields(module).Errorf(
			"%v: Polling interval should be > 0 second. Given PollingIntervalInSeconds: %v",
			name, pollingIntervalInSeconds)
		return nil
	}

	log.WithFields(module).Infof(
		"Creating Global Dispatcher - %v. MaxWorkers: %v | PollingIntervalInSeconds: %v",
		name, maxWorkers, pollingIntervalInSeconds)

	workerPool := make(chan chan Job, maxWorkers)

	return &GlobalDispatcher{
		Name:                     name,
		JobConsumer:              jobConsumer,
		WorkerPool:               workerPool,
		MaxWorkers:               maxWorkers,
		PollingIntervalInSeconds: pollingIntervalInSeconds,
		shutdown:                 make(chan bool, 1),
		confirmShutdown:          make(chan bool, 1),
	}
}

// GetName - Gets the name of the dispatcher
func (d *GlobalDispatcher) GetName() string {
	return d.Name
}

// Run - Dispatches the workers as per the configuration
func (d *GlobalDispatcher) Run() {
	run(d, d.MaxWorkers, d.WorkerPool)
}

// Shutdown - Gracefully shuts down the dispatcher
// by closing the job enqueuing channel
// and waits for workers to complete the accepted jobs
func (d *GlobalDispatcher) Shutdown() <-chan struct{} {
	log.WithFields(module).Infof(
		"%v: Shutting down the Global Dispatcher", d.Name)

	// Sending the signal to shutdown the dispatcher
	d.shutdown <- true

	shutdownComplete := make(chan struct{})

	// Wait for all the pending jobs in the queue to get served
	<-d.confirmShutdown

	// No jobs in the queue. Stop the workers and
	// wait for their graceful shutdown
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
		"%v: Gracefully shut down the Global Dispatcher", d.Name)

	return shutdownComplete
}

func (d *GlobalDispatcher) addWorker(w *Worker) {
	d.workers = append(d.workers, w)
}

func (d *GlobalDispatcher) dispatch() {
	for {
		select {
		case jobChannel := <-d.WorkerPool:
			// Worker is ready
			// No job is fetched for processing till the worker is ready
			for {
				// Get the job using the provided custom job consumer object
				job := d.JobConsumer.Consume()
				if job != nil {
					log.WithFields(module).Debugf("%v: Received a job with job ID %v", d.Name, job.GetJobID())

					// Add the job to the corresponding worker's job channel
					jobChannel <- job
					break
				} else {
					// No jobs to work on
					select {
					case <-d.shutdown:
						// We have received a signal to stop
						// Confirm the dispatcher shutdown
						d.confirmShutdown <- true
						return
					default:
						// Poll after the specified polling interval
						log.WithFields(module).Debugf("%v: No job received. Will  poll the queue after %v seconds",
							d.Name, d.PollingIntervalInSeconds)

						time.Sleep(
							time.Duration(
								d.PollingIntervalInSeconds) * 1000 * time.Millisecond)
					}
				}
			}
		}
	}
}
