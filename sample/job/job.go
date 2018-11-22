package job

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var module = log.Fields{"module": "my_job"}

// MyJob ...
type MyJob struct {
	JobID string
	X     int
	Y     int
}

// GetJobID ...
func (j *MyJob) GetJobID() string {
	return j.JobID
}

// GetTimeout ...
func (j *MyJob) GetTimeout() time.Duration {
	return 1 * time.Second
}

// DoJob ...
func (j *MyJob) DoJob() []error {
	log.WithFields(module).Infof("My Job %v Started", j.JobID)

	time.Sleep(1000 * time.Millisecond)
	log.WithFields(module).Infof("X[%v]+Y[%v]=%v", j.X, j.Y, (j.X + j.Y))

	log.WithFields(module).Infof("My Job %v Completed", j.JobID)

	return nil
}
