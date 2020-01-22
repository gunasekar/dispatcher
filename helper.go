package dispatcher

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var module = log.Fields{"module": "dispatcher"}

func newUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", 4)
}

func run(d Dispatcher, maxWorkers int, workerPool chan chan Job) {
	if d == nil {
		log.WithFields(module).Errorf("Dispatcher object is null")
		panic("Dispatcher object is null")
	}

	if workerPool == nil {
		log.WithFields(module).Errorf("WorkerPool is null")
		panic("WorkerPool is null")
	}

	log.WithFields(module).Infof("Creating %v Workers for the dispatcher - %v",
		maxWorkers,
		d.GetName())

	for i := 0; i < maxWorkers; i++ {
		workerID := d.GetName() + "-Worker-" + strconv.Itoa(i)
		worker := NewWorker(workerID, workerPool)
		d.addWorker(worker)
		worker.Start()
	}

	go d.dispatch()
}
