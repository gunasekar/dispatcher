package dispatcher

import (
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test_HelperModule_run1(t *testing.T) {
	convey.Convey("Call the run helper function with null dispatcher object", t, func() {
		convey.So(func() { run(nil, 2, make(chan chan Job, 2)) }, convey.ShouldPanic)
	})
}

func Test_HelperModule_run2(t *testing.T) {
	jobQueue := make(chan Job, 2)
	dispatcher := NewLocalDispatcher("testLD", 1, jobQueue)
	convey.Convey("Call the run helper function with null worker pool", t, func() {
		convey.So(func() { run(dispatcher, 2, nil) }, convey.ShouldPanic)
	})
}
