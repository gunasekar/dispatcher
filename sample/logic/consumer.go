package logic

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"gitlab.com/dispatcher"
)

// TestJobConsumer ...
type TestJobConsumer struct {
	Values chan int
}

// Consume ...
func (jc *TestJobConsumer) Consume() dispatcher.Job {
	result, err := SqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &QueueURL,
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

	job := &MyJob{}
	json.Unmarshal([]byte(*result.Messages[0].Body), job)
	return job
}
