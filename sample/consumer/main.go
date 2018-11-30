package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gunasekar/dispatcher"
	"github.com/gunasekar/dispatcher/sample/consumer/awssqs"
	"github.com/gunasekar/dispatcher/sample/consumer/awssqs/deleter"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Create SQS client
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("foo", "bar", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.UsWest2RegionID),
		Endpoint:         aws.String("http://localhost:4576"),
	}))

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	sqsClient := sqs.New(sess)
	queueURL := "http://localhost:4576/queue/localq"

	// Create the object which defines to the consume logic
	var jobConsumer dispatcher.JobConsumer
	jobConsumer = &awssqs.TestJobConsumer{SqsClient: sqsClient,
		QueueURL: queueURL}

	deleter.SQSJobDeleter = &deleter.JobDeleter{SqsClient: sqsClient,
		QueueURL: queueURL}

	// Create and start the dispatcher
	dispatcher := dispatcher.NewGlobalDispatcher("test-dispatcher", 3, jobConsumer, 1)
	dispatcher.Run()

	// wait for syscall.SIGINT or syscall.SIGTERM and shutdown the dispatcher
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)
	<-interrupt
	<-dispatcher.Shutdown()

}
