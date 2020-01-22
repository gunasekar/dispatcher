package dispatcher

import (
	"time"

	"github.com/google/uuid"
)

// InjectedJob defines a generic job where any action can be injected
type InjectedJob struct {
	JobID                            string
	Action                           func()
	FinallyAction                    func()
	ExecutionTimeoutInSeconds        int64
	FinallyExecutionTimeoutInSeconds int64
}

// NewJob instantiates new job object
func NewJob(action func(), execTimeoutInSec, finallyTimeoutInSec int64) *InjectedJob {
	return &InjectedJob{
		JobID:                            uuid.New().String(),
		Action:                           action,
		ExecutionTimeoutInSeconds:        execTimeoutInSec,
		FinallyExecutionTimeoutInSeconds: finallyTimeoutInSec,
	}
}

// NewJobWithID instantiates new job object with jobID
func NewJobWithID(jobID string, action func(), execTimeoutInSec, finallyTimeoutInSec int64) *InjectedJob {
	return &InjectedJob{
		JobID:                            jobID,
		Action:                           action,
		ExecutionTimeoutInSeconds:        execTimeoutInSec,
		FinallyExecutionTimeoutInSeconds: finallyTimeoutInSec,
	}
}

// GetJobID - Returns JobID
func (j *InjectedJob) GetJobID() string {
	return j.JobID
}

// Execute executes the injected action
func (j *InjectedJob) Execute() []error {
	j.Action()
	return nil
}

// GetExecutionTimeout get the execution timeout of the action
func (j *InjectedJob) GetExecutionTimeout() time.Duration {
	return time.Duration(j.ExecutionTimeoutInSeconds) * time.Second
}

// Finally executes at the end of the action
func (j *InjectedJob) Finally() {
	if j.FinallyAction != nil {
		j.FinallyAction()
	}
}

// GetFinallyTimeout get the finally timeout of the action
func (j *InjectedJob) GetFinallyTimeout() time.Duration {
	return time.Duration(j.FinallyExecutionTimeoutInSeconds) * time.Second
}
