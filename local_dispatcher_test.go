package dispatcher

import (
	"errors"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
)

type TestJob struct {
	JobID string
}

func (j *TestJob) GetJobID() string {
	return j.JobID
}

func (j *TestJob) Execute() []error {
	var errs []error

	if j.JobID == "-1" {
		log.WithFields(module).Debugf("Creating a NULL ref exception")
		var err1 error
		_ = err1.Error()
	}

	if j.JobID == "" {
		log.WithFields(module).Debugf("Error in JOB")
		errs = append(errs, errors.New("Nil Job ID"))
		return errs
	}

	log.WithFields(module).Debugf("Test Job %v Started", j.JobID)
	time.Sleep(2 * time.Second)
	log.WithFields(module).Debugf("Test Job %v Completed", j.JobID)

	return nil
}

func (j *TestJob) GetExecutionTimeout() time.Duration {
	return 10 * time.Second
}

func (j *TestJob) Finally() {
	log.WithFields(module).Debugf("Running finally")
}

func (j *TestJob) GetFinallyTimeout() time.Duration {
	return 1 * time.Second
}

func Test_LocalDispatcher_SunnyDay_InjectedJob(t *testing.T) {
	convey.Convey("Testing Local dispatcher's sunny day scenario", t, func() {
		convey.So(func() {
			jobQueue := make(chan Job, 10)
			dispatcher := NewLocalDispatcher("testLD", 1, jobQueue)
			dispatcher.Run()

			for i := 0; i < 3; i++ {
				jobQueue <- NewJob(func() {
					log.WithFields(module).Debugf("Injected Job %v Started", i)
					time.Sleep(2 * time.Second)
					log.WithFields(module).Debugf("Injected Job %v Completed", i)
				}, 10, 1)
				log.WithFields(module).Debugf("Queued Job %v", i)
				time.Sleep(1 * time.Second)
			}

			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}

func Test_LocalDispatcher_SunnyDay(t *testing.T) {
	convey.Convey("Testing Local dispatcher's sunny day scenario", t, func() {
		convey.So(func() {
			jobQueue := make(chan Job, 10)
			dispatcher := NewLocalDispatcher("testLD", 1, jobQueue)
			dispatcher.Run()

			for i := 0; i < 3; i++ {
				jobQueue <- &TestJob{JobID: strconv.Itoa(i)}
				log.WithFields(module).Debugf("Queued Job %v", i)
				time.Sleep(1 * time.Second)
			}

			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}

func Test_LocalDispatcher_WithJobReturningErrorf(t *testing.T) {
	convey.Convey("Testing Local dispatcher with job returning error", t, func() {
		convey.So(func() {
			jobQueue := make(chan Job, 2)
			dispatcher := NewLocalDispatcher("", 1, jobQueue)
			dispatcher.Run()
			jobQueue <- &TestJob{JobID: ""}
			log.WithFields(module).Debugf("Queued Job with empty ID")
			time.Sleep(3 * time.Second)
			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}

func Test_LocalDispatcher_WithNoJobChannel(t *testing.T) {
	convey.Convey("Local Dispatcher with no job channel", t, func() {
		dispatcher := NewLocalDispatcher("testLD", 2, nil)
		convey.So(dispatcher, convey.ShouldBeNil)
	})
}

func Test_LocalDispatcher_WithNoWorkers(t *testing.T) {
	convey.Convey("Local Dispatcher with no workers", t, func() {
		dispatcher := NewLocalDispatcher("testLD", 0, nil)
		convey.So(dispatcher, convey.ShouldBeNil)
	})
}

func Test_LocalDispatcher_JobThrowingException(t *testing.T) {
	convey.Convey("Testing Local Dispatcher with job throwing exception", t, func() {
		convey.So(func() {
			jobQueue := make(chan Job, 2)
			dispatcher := NewLocalDispatcher("testLD", 1, jobQueue)
			dispatcher.Run()
			jobQueue <- &TestJob{JobID: "-1"}
			time.Sleep(3 * time.Second)
			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}
