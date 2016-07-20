package object

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	platformSQS "github.com/tapglue/multiverse/platform/sqs"
)

const (
	queueName = "object-state-change"
)

type sqsSource struct {
	api      platformSQS.API
	queueURL string
}

// SQSSource returns an SQS backed Source implementation.
func SQSSource(api platformSQS.API) (Source, error) {
	res, err := api.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}

	return &sqsSource{
		api:      api,
		queueURL: *res.QueueUrl,
	}, nil
}

func (s *sqsSource) Ack(id string) error {
	_, err := s.api.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: aws.String(id),
	})

	return err
}

func (s *sqsSource) Consume() (*StateChange, error) {
	o, err := s.api.ReceiveMessage(&sqs.ReceiveMessageInput{
		MessageAttributeNames: []*string{
			aws.String(platformSQS.AttributeAll),
		},
		QueueUrl:          aws.String(s.queueURL),
		VisibilityTimeout: aws.Int64(platformSQS.TimeoutVisibility),
		WaitTimeSeconds:   aws.Int64(platformSQS.TimeoutWait),
	})
	if err != nil {
		return nil, err
	}

	if len(o.Messages) == 0 {
		return nil, ErrEmptySource
	}

	var (
		m = o.Messages[0]

		sentAt time.Time
	)

	if attr, ok := m.MessageAttributes[platformSQS.AttributeSentAt]; ok {
		t, err := time.Parse(platformSQS.FormatSentAt, *attr.StringValue)
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
		AckID:     *m.ReceiptHandle,
		ID:        *m.MessageId,
		Namespace: f.Namespace,
		New:       f.New,
		Old:       f.Old,
		SentAt:    sentAt,
	}, nil
}

func (s *sqsSource) Propagate(ns string, old, new *Object) (string, error) {
	r, err := json.Marshal(&stateChange{
		Namespace: ns,
		New:       new,
		Old:       old,
	})
	if err != nil {
		return "", err
	}

	o, err := s.api.SendMessage(s.messageInput(r))
	if err != nil {
		return "", err
	}

	return *o.MessageId, nil
}

func (s *sqsSource) messageInput(body []byte) *sqs.SendMessageInput {
	now := time.Now().Format(platformSQS.FormatSentAt)

	return &sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			platformSQS.AttributeSentAt: &sqs.MessageAttributeValue{
				DataType:    aws.String(platformSQS.TypeString),
				StringValue: aws.String(now),
			},
		},
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(s.queueURL),
	}
}

type stateChange struct {
	Namespace string  `json:"namespace"`
	New       *Object `json:"new"`
	Old       *Object `json:"old"`
}
