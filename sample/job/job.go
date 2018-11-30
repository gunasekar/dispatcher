package job

import (
	"time"

	"github.com/gunasekar/dispatcher/sample/consumer/awssqs/deleter"
	log "github.com/sirupsen/logrus"
)

var module = log.Fields{"module": "my_job"}

// MyJob ...
type MyJob struct {
	JobID         string
	X             int
	Y             int
	ReceiptHandle string
}

// GetJobID ...
func (j *MyJob) GetJobID() string {
	return j.JobID
}

// Execute ...
func (j *MyJob) Execute() []error {
	log.WithFields(module).Infof("My Job %v Started", j.JobID)

	time.Sleep(100 * time.Millisecond)
	log.WithFields(module).Infof("X[%v]+Y[%v]=%v", j.X, j.Y, (j.X + j.Y))

	log.WithFields(module).Infof("My Job %v Completed", j.JobID)

	return nil
}

// GetExecutionTimeout ...
func (j *MyJob) GetExecutionTimeout() time.Duration {
	return 1 * time.Second
}

// Finally ...
func (j *MyJob) Finally() {
	log.WithFields(module).Debugf("Running finally")
	deleter.SQSJobDeleter.DeleteMessage(j.ReceiptHandle)
	//time.Sleep(100 * time.Millisecond)
}

// GetFinallyTimeout ...
func (j *MyJob) GetFinallyTimeout() time.Duration {
	return 1 * time.Second
}
