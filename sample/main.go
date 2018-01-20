package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/dispatcher"
	"gitlab.com/dispatcher/sample/logic"
)

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	logic.SqsClient = sqs.New(sess)
	logic.QueueURL = "https://sqs.ap-southeast-1.amazonaws.com/903916081954/test-q"
}

func main() {
	var jobConsumer dispatcher.JobConsumer
	jobConsumer = &logic.TestJobConsumer{}
	dispatcher := dispatcher.NewGlobalDispatcher("test-dispatcher", 5, jobConsumer, 3)
	dispatcher.Run()

	startProducing()

	<-dispatcher.Shutdown()
}

func startProducing() {
	go func() {
		for i := 0; i < 5; i++ {
			job := &logic.MyJob{JobID: uuid.NewV4().String(), X: i + 1, Y: i + 2}
			jobContent, _ := json.Marshal(job)
			logic.SqsClient.SendMessage(&sqs.SendMessageInput{
				DelaySeconds: aws.Int64(10),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"jobID": &sqs.MessageAttributeValue{
						DataType:    aws.String("String"),
						StringValue: aws.String(job.GetJobID()),
					},
				},
				MessageBody: aws.String(string(jobContent)),
				QueueUrl:    &logic.QueueURL,
			})

			log.Debugf("Pushed a job to queue: %#v", job)
			time.Sleep(5 * time.Second)
		}
	}()
}
