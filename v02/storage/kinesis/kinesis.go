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
		// SetupStreams creates the needed Kinesis streams
		SetupStreams([]string) error

		// PutRecord sends a new record to a Kinesis stream
		PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, error)

		// TeardownStreams destroys the streams from Kinesis
		TeardownStreams(streamsName []string) error

		// Datastore returns the Kinesis client
		Datastore() *gksis.Kinesis
	}

	cli struct {
		kinesis *gksis.Kinesis
	}
)

const (
	// StreamNewAccount stream name
	StreamNewAccount = "new_account"
)

var (
	// Streams defines the array of all streams currently defined by the application
	Streams = []string{StreamNewAccount}
)

func (c *cli) PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, error) {
	time.Sleep(time.Duration(rand.Intn(100)*100) * time.Millisecond)
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.AddRecord(payload, partitionKey)
	return c.kinesis.PutRecord(args)
}

func (c *cli) SetupStreams(streamsName []string) error {
	shardCount := 1 // TODO this should be configurable maybe?
	for _, streamName := range streamsName {
		err := c.kinesis.CreateStream(streamName, shardCount)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (c *cli) TeardownStreams(streamsName []string) error {
	for _, streamName := range streamsName {
		err := c.kinesis.DeleteStream(streamName)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (c *cli) Datastore() *gksis.Kinesis {
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

// NewTest returns a new testing-enabled client
func NewTest(authKey, secretKey, region, endpoint string) Client {
	auth := &gksis.Auth{
		AccessKey: authKey,
		SecretKey: secretKey,
	}

	return &cli{
		kinesis: gksis.NewWithEndpoint(auth, gksis.Region{Name: region}, endpoint),
	}
}
