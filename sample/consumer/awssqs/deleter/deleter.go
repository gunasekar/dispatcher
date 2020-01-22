package deleter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("Delete message failed for receiptHandle - " + receiptHandle)
	}

	log.Debugf("DeletedOutput: %v", deleteMsgOutput.GoString())
}
