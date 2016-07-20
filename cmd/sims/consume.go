package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
)

const typeNotification = "Notification"

type sqsReceiver interface {
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
}

type conRuleFunc func(*connection.StateChange) (*message, error)
type objectRuleFunc func(*object.StateChange) ([]*message, error)

type endpointChange struct {
	ack            ackFunc
	EndpointArn    string `json:"EndpointArn"`
	EventType      string `json:"EventType"`
	FailureMessage string `json:"FailureMessage"`
	FailureType    string `json:"FailureType"`
	Resource       string `json:"Resource"`
	Service        string `json:"Service"`
}

func conRuleFollower(
	fetchUser fetchUserFunc,
) conRuleFunc {
	return func(change *connection.StateChange) (*message, error) {
		if change.Old != nil ||
			change.New.State != connection.StateConfirmed ||
			change.New.Type != connection.TypeFollow {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.FromID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		target, err := fetchUser(change.Namespace, change.New.ToID)
		if err != nil {
			return nil, fmt.Errorf("target fetch: %s", err)
		}

		return &message{
			message: fmt.Sprintf(
				"%s %s (%s) started following you",
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			recipient: target.ID,
		}, nil
	}
}

func consumeConnection(
	conSource connection.Source,
	batchc chan<- batch,
	ruleFns ...conRuleFunc,
) error {
	for {
		c, err := conSource.Consume()
		if err != nil {
			if connection.IsEmptySource(err) {
				continue
			}
			return err
		}

		ms := []*message{}

		for _, rule := range ruleFns {
			msg, err := rule(c)
			if err != nil {
				return err
			}

			if msg != nil {
				ms = append(ms, msg)
			}
		}

		if len(ms) == 0 {
			err := conSource.Ack(c.AckID)
			if err != nil {
				return err
			}

			continue
		}

		batchc <- batch{
			ackFunc: func() error {
				acked := false

				if acked {
					return nil
				}

				err := conSource.Ack(c.AckID)
				if err == nil {
					acked = true
				}
				return err
			},
			messages:  ms,
			namespace: c.Namespace,
		}
	}
}

func consumeEndpointChange(r sqsReceiver, queueURL string, changec chan endpointChange) error {
	for {
		o, err := r.ReceiveMessage(&sqs.ReceiveMessageInput{
			MessageAttributeNames: []*string{
				aws.String("All"),
			},
			QueueUrl:          aws.String(queueURL),
			VisibilityTimeout: aws.Int64(5),
			WaitTimeSeconds:   aws.Int64(10),
		})
		if err != nil {
			return err
		}

		for _, msg := range o.Messages {
			ack := func() error {
				_, err := r.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				return err
			}

			if msg.Body == nil {
				_ = ack()

				continue
			}

			f := struct {
				Message string `json:"Message"`
				Type    string `json:"Type"`
			}{}

			if err := json.Unmarshal([]byte(*msg.Body), &f); err != nil {
				return err
			}

			if f.Type != typeNotification {
				_ = ack()

				continue
			}

			c := endpointChange{}

			if err := json.Unmarshal([]byte(f.Message), &c); err != nil {
				return err
			}

			c.ack = ack

			changec <- c
		}

	}
}

func consumeObject(
	objectSource object.Source,
	batchc chan<- batch,
	ruleFns ...objectRuleFunc,
) error {
	for {
		c, err := objectSource.Consume()
		if err != nil {
			if object.IsEmptySource(err) {
				continue
			}
			return err
		}

		ms := []*message{}

		for _, rule := range ruleFns {
			rs, err := rule(c)
			if err != nil {
				return err
			}

			for _, msg := range rs {
				ms = append(ms, msg)
			}
		}

		if len(ms) == 0 {
			err = objectSource.Ack(c.AckID)
			if err != nil {
				return err
			}

			continue
		}

		batchc <- batch{
			ackFunc: func() error {
				acked := false

				if acked {
					return nil
				}

				err = objectSource.Ack(c.AckID)
				if err == nil {
					acked = true
				}
				return err
			},
			messages:  ms,
			namespace: c.Namespace,
		}

		err = objectSource.Ack(c.AckID)
		if err != nil {
			return err
		}
	}
}

func objectRuleCommentCreated(
	fetchFriends fetchFriendsFunc,
	fetchUser fetchUserFunc,
) objectRuleFunc {
	return func(change *object.StateChange) ([]*message, error) {
		if change.Old != nil ||
			change.New.Deleted == true ||
			!isComment(change.New) {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		fs, err := fetchFriends(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, friend := range fs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					"Your friend %s %s (%s) commented on a Post.",
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: friend.ID,
			})
		}

		return ms, nil
	}
}

func objectRulePostCreated(
	fetchFriends fetchFriendsFunc,
	fetchUser fetchUserFunc,
) objectRuleFunc {
	return func(change *object.StateChange) ([]*message, error) {
		if change.Old != nil ||
			!isPost(change.New) ||
			change.New.Deleted == true {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		fs, err := fetchFriends(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, friend := range fs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					"Your friend %s %s (%s) created a new Post.",
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: friend.ID,
			})
		}

		return ms, nil
	}
}

func isComment(o *object.Object) bool {
	if o.Type != controller.TypeComment {
		return false
	}

	return o.Owned
}

func isPost(o *object.Object) bool {
	if o.Type != controller.TypePost {
		return false
	}

	return o.Owned
}
