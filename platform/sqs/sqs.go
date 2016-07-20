package sqs

import "github.com/aws/aws-sdk-go/service/sqs"

// Common Attributes.
const (
	AttributeSentAt = "SentAt"
	AttributeAll    = "All"

	FormatSentAt = "2006-01-02 15:04:05.999999999 -0700 MST"

	TypeString = "String"
)

// Common Timeouts.
var (
	TimeoutVisibility int64 = 60
	TimeoutWait       int64 = 10
)

// API bundles common SQS operations.
type API interface {
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}
