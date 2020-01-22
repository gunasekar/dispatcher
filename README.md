# dispatcher [![Build Status](https://travis-ci.org/gunasekar/dispatcher.svg?branch=master)](https://travis-ci.org/gunasekar/dispatcher)

golang based dispatcher library inspired by http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

## local dispatcher
managed worker execution with GO channel as job source

```
// initialize job queue to hold max 100 jobs. basically trying to restrict the job producers
jobQueue := make(chan Job, 100)

// create the local dispatcher with 10 workers
jobDispatcher := NewLocalDispatcher("my_local_dispatcher", 10, jobQueue)

// start dispatching
jobDispatcher.Run()

// produce the jobs to execute in a managed fashion
for i := 0; i < 3; i++ {
  jobQueue <- NewJob(func() {
    // action to be delegated to the worker
  }, 10, 1)
}

// shutdown the dispatcher
<-jobDispatcher.Shutdown()
 ```

## global dispatcher
managed worker execution with any external job source


```
// define your job consumer
type TestJobConsumer struct {
	Values chan int
}

func (jc *TestJobConsumer) Consume() Job {
	if jc.Values == nil {
		return nil
	}

	select {
	case i := <-jc.Values:
		return NewJob(func() {
			fmt.Printf("Executing %v", i)
            // action to be delegated to the worker
		}, 10, 1)
		//&TestJob{JobID: strconv.Itoa(i)}
	default:
		return nil
	}
}
```

```
// initialize job queue
jobQueue := make(chan Job)

// instantiate a custom job consumer
jc := &TestJobConsumer{Values: make(chan int, 10)}

// create the global dispatcher with the above created job consumer
// assigning 10 workers and pollingIntervalInSeconds of 2s when no job found
jobDispatcher := NewGlobalDispatcher("my_global_dispatcher", 10, jc, 2)

// start dispatching
jobDispatcher.Run()

// shutdown the dispatcher
<-jobDispatcher.Shutdown()
 ```
