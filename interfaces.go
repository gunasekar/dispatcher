package dispatcher

import "time"

// Dispatcher ...
type Dispatcher interface {
	GetName() string
	Run()
	Shutdown() <-chan struct{}
	dispatch()
	addWorker(*Worker)
}

// Job ...
type Job interface {
	GetJobID() string
	Execute() []error
	GetExecutionTimeout() time.Duration
	Finally()
	GetFinallyTimeout() time.Duration
}

// JobConsumer ...
type JobConsumer interface {
	Consume() Job
}
