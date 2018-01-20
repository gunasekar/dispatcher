package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gitlab.com/dispatcher"
	"gitlab.com/dispatcher/sample/consumer/awssqs"
)

func main() {
	// Create SQS client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sess)
	queueURL := "https://sqs.ap-southeast-1.amazonaws.com/903916081954/test-q"

	// Create the object which defines to the consume logic
	var jobConsumer dispatcher.JobConsumer
	jobConsumer = &awssqs.TestJobConsumer{SqsClient: sqsClient,
		QueueURL: queueURL}

	// Create and start the dispatcher
	dispatcher := dispatcher.NewGlobalDispatcher("test-dispatcher", 5, jobConsumer, 3)
	dispatcher.Run()

	// wait for syscall.SIGINT or syscall.SIGTERM and shutdown the dispatcher
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)
	<-interrupt
	<-dispatcher.Shutdown()

}
