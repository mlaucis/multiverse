/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package kinesis provides the AWS Kinesis needed functions for kinesis
package kinesis

import (
	"math/rand"
	"time"

	gksis "github.com/sendgridlabs/go-kinesis"
)

type (
	// Client defines the interface for Kinesis
	Client interface {
		Client() *gksis.Kinesis
		SetupStreams([]string) error
		PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, error)
	}

	cli struct {
		kinesis *gksis.Kinesis
	}
)

const (
	// StreamNewAccount stream name
	StreamNewAccount = "new_account"
)

// PutRecord sends a new record to a Kinesis stream
func (c *cli) PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, error) {
	time.Sleep(time.Duration(rand.Intn(100)*100) * time.Millisecond)
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.AddRecord(payload, partitionKey)
	return c.kinesis.PutRecord(args)
}

// SetupStreams creates the needed Kinesis streams
func (c *cli) SetupStreams(streamsName []string) error {
	for _, streamName := range streamsName {
		err := c.kinesis.CreateStream(streamName, 1)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// Client returns the Kinesis client
func (c *cli) Client() *gksis.Kinesis {
	return c.kinesis
}

// New returns a new Kinesis client
func New(authKey, secretKey, region string) Client {
	auth := &gksis.Auth{
		AccessKey: authKey,
		SecretKey: secretKey,
	}

	return &cli{
		kinesis: gksis.New(auth, gksis.Region{Name: region}),
	}
}
