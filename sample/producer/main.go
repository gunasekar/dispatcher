package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/dispatcher/sample/job"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Configure Queue
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqsClient := sqs.New(sess)
	queueURL := "https://sqs.ap-southeast-1.amazonaws.com/903916081954/test-q"

	// Start producing messages
	for i := 0; i < 50; i++ {
		job := &job.MyJob{JobID: uuid.NewV4().String(), X: i + 1, Y: i + 2}
		jobContent, _ := json.Marshal(job)
		sqsClient.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"jobID": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(job.GetJobID()),
				},
			},
			MessageBody: aws.String(string(jobContent)),
			QueueUrl:    &queueURL,
		})

		log.Debugf("Pushed a job to queue: %#v", job)
		time.Sleep(2 * time.Second)
	}
}
