package dispatcher

import (
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
)

type TestJobConsumer struct {
	Values chan int
}

func (jc *TestJobConsumer) Consume() Job {
	if jc.Values == nil {
		return nil
	}

	select {
	case i := <-jc.Values:
		return &TestJob{JobID: strconv.Itoa(i)}
	default:
		return nil
	}
}

func Test_GlobalDispatcher_SunnyDay(t *testing.T) {
	convey.Convey("Testing Global dispatcher's sunny day scenario", t, func() {
		convey.So(func() {
			jc := &TestJobConsumer{Values: make(chan int, 10)}
			dispatcher := NewGlobalDispatcher("testGD", 1, jc, 1)
			dispatcher.Run()
			for i := 0; i < 3; i++ {
				jc.Values <- i
				log.WithFields(module).Debugf("Queued Job %v", i)
				time.Sleep(1 * time.Second)
			}

			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}

func Test_GlobalDispatcher_WithNoJobConsumer(t *testing.T) {
	convey.Convey("Dispatcher with no job consumer", t, func() {
		dispatcher := NewGlobalDispatcher("", 1, nil, 1)
		convey.So(dispatcher, convey.ShouldBeNil)
	})
}

func Test_GlobalDispatcherWithNoWorkers(t *testing.T) {
	convey.Convey("Dispatcher creation with 0 workers", t, func() {
		dispatcher := NewGlobalDispatcher("testGD", 0, nil, 1)
		convey.So(dispatcher, convey.ShouldBeNil)
	})
}

func Test_GlobalDispatcherWith0SecondPolling(t *testing.T) {
	convey.Convey("Dispatcher creation with no polling interval", t, func() {
		jc := &TestJobConsumer{}
		dispatcher := NewGlobalDispatcher("testGD", 1, jc, 0)
		convey.So(dispatcher, convey.ShouldBeNil)
	})
}

func Test_GlobalDispatcherWithNoJobHencePoll(t *testing.T) {
	convey.Convey("Testing polling when dispatcher finds no job", t, func() {
		convey.So(func() {
			jc := &TestJobConsumer{}
			dispatcher := NewGlobalDispatcher("testGD", 1, jc, 1)
			dispatcher.Run()
			time.Sleep(3 * time.Second)
			<-dispatcher.Shutdown()
		}, convey.ShouldNotPanic)
	})
}
