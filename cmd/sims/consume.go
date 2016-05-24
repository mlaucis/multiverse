package main

import (
	"fmt"

	"github.com/tapglue/multiverse/service/connection"
)

type conRuleFunc func(*connection.StateChange) (*message, error)

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
				err := conSource.Ack(c.AckID)
				if err != nil {
					return err
				}

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
