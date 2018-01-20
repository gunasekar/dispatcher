package awssqs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"gitlab.com/dispatcher"
	"gitlab.com/dispatcher/sample/job"
)

// TestJobConsumer ...
type TestJobConsumer struct {
	SqsClient *sqs.SQS
	QueueURL  string
}

// Consume ...
func (jc *TestJobConsumer) Consume() dispatcher.Job {
	result, err := jc.SqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &jc.QueueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(36000),
		WaitTimeSeconds:     aws.Int64(20),
	})

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil
	}

	if len(result.Messages) == 0 {
		log.Debugf("Received no messages")
		return nil
	}

	job := &job.MyJob{}
	json.Unmarshal([]byte(*result.Messages[0].Body), job)
	return job
}
