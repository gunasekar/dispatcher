package deleter

import (
	"context"

	"bitbucket.org/dreamplug-backend/commons-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSJobDeleter ...
var SQSJobDeleter *JobDeleter

// JobDeleter ...
type JobDeleter struct {
	SqsClient *sqs.SQS
	QueueURL  string
}

// DeleteMessage ...
func (d *JobDeleter) DeleteMessage(receiptHandle string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(d.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	deleteMsgOutput, err := d.SqsClient.DeleteMessage(params)
	if err != nil {
		logger.Errorf(context.Background(), "Delete message failed for receiptHandle - "+receiptHandle)
	}

	logger.Debugf(context.Background(), "DeletedOutput: %v", deleteMsgOutput.GoString())
}
