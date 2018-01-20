package dispatcher

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_Worker_nil_WorkerPool(t *testing.T) {
	convey.Convey("Worker without worker pool", t, func() {
		w := NewWorker("Worker_1", nil)
		convey.So(w, convey.ShouldBeNil)
	})
}

func Test_Worker_no_worker_id(t *testing.T) {
	convey.Convey("Worker without worker id", t, func() {
		w := NewWorker("", make(chan chan Job, 1))
		convey.So(w, convey.ShouldBeNil)
	})
}
