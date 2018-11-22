package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gunasekar/dispatcher/sample/job"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Configure Queue
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("foo", "var", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.UsWest2RegionID),
		Endpoint:         aws.String("http://localhost:4576"),
	}))

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	sqsClient := sqs.New(sess)
	queueURL := "http://localhost:4576/queue/localq"

	// Start producing messages
	for i := 0; i < 5; i++ {
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
