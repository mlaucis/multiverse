package connection

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	attributeSentAt  = "SentAt"
	sentAtFormat     = "2006-01-02 15:04:05.999999999 -0700 MST"
	sqsAttribtuesAll = "All"
	sqsTypeString    = "String"
)

var (
	visibilityTimeout int64 = 1
	waitTimeSeconds   int64 = 1
)

type sendReceiver interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type sqsSource struct {
	api      sendReceiver
	queueURL string
}

// SQSSource reutrns an SQS backed Source implementation.
func SQSSource(api sendReceiver, queueURL string) Source {
	return &sqsSource{
		api:      api,
		queueURL: queueURL,
	}
}

func (s *sqsSource) Consume() (*StateChange, error) {
	all := sqsAttribtuesAll

	o, err := s.api.ReceiveMessage(&sqs.ReceiveMessageInput{
		MessageAttributeNames: []*string{
			&all,
		},
		QueueUrl:          &s.queueURL,
		VisibilityTimeout: &visibilityTimeout,
		WaitTimeSeconds:   &waitTimeSeconds,
	})
	if err != nil {
		return nil, err
	}

	var (
		m = o.Messages[0]

		sentAt time.Time
	)

	if attr, ok := m.MessageAttributes[attributeSentAt]; ok {
		t, err := time.Parse(sentAtFormat, *attr.StringValue)
		if err != nil {
			return nil, err
		}

		sentAt = t
	}

	f := stateChange{}

	err = json.Unmarshal([]byte(*m.Body), &f)
	if err != nil {
		return nil, err
	}

	return &StateChange{
		ID:        *m.MessageId,
		Namespace: f.Namespace,
		New:       f.New,
		Old:       f.Old,
		SentAt:    sentAt,
	}, nil
}

func (s *sqsSource) Propagate(ns string, old, new *Connection) (string, error) {
	r, err := json.Marshal(&stateChange{
		Namespace: ns,
		New:       new,
		Old:       old,
	})
	if err != nil {
		return "", err
	}

	o, err := s.api.SendMessage(s.messageInput(string(r)))
	if err != nil {
		return "", err
	}

	return *o.MessageId, nil
}

func (s *sqsSource) messageInput(body string) *sqs.SendMessageInput {
	var (
		now        = time.Now().Format(sentAtFormat)
		typeString = sqsTypeString
	)

	return &sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			attributeSentAt: &sqs.MessageAttributeValue{
				DataType:    &typeString,
				StringValue: &now,
			},
		},
		MessageBody: &body,
		QueueUrl:    &s.queueURL,
	}
}

type stateChange struct {
	Namespace string      `json:"namespace"`
	New       *Connection `json:"new"`
	Old       *Connection `json:"old"`
}
