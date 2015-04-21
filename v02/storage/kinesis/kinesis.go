/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package kinesis provides the AWS Kinesis needed functions for kinesis
package kinesis

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tapglue/backend/tgerrors"

	gksis "github.com/sendgridlabs/go-kinesis"
)

type (
	// Client defines the interface for Kinesis
	Client interface {
		// SetupStreams creates the needed Kinesis streams
		SetupStreams([]string) error

		// PutRecord sends a new record to a Kinesis stream
		PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, tgerrors.TGError)

		// GetRecords returns at most the specified number of records from the desired stream / shard
		GetRecords(streamName, shardID, consumerName string, maxEntries int) ([]string, tgerrors.TGError)

		// StreamRecords will stream all the records it can from from all the shards of a stream
		// To stop it, just close the output channel
		StreamRecords(streamName, consumerName string, output <-chan string, errors <-chan tgerrors.TGError)

		// DescribeStream will return the Kinesis stream descriptiont that AWS has if the stream is active.
		// If the stream is not active then it will return an error
		// It will also timeout after 30 seconds (one try per second)
		DescribeStream(streamName string) (*gksis.DescribeStreamResp, tgerrors.TGError)

		// TeardownStreams destroys the streams from Kinesis
		TeardownStreams(streamsName []string) error

		// Datastore returns the Kinesis client
		Datastore() *gksis.Kinesis
	}

	cli struct {
		kinesis *gksis.Kinesis
	}
)

func (c *cli) PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, tgerrors.TGError) {
	time.Sleep(time.Duration(rand.Intn(100)*100) * time.Millisecond)
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.AddRecord(payload, partitionKey)
	resp, err := c.kinesis.PutRecord(args)
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to execute operation (1)", err.Error())
	}

	return resp, nil
}

func (c *cli) GetRecords(streamName, shardID, consumerName string, maxEntries int) ([]string, tgerrors.TGError) {
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", shardID)
	args.Add("ShardIteratorType", "TRIM_HORIZON")
	shardIteratorResponse, err := c.kinesis.GetShardIterator(args)
	if err != nil {
		return "", tgerrors.NewInternalError("error while reading the internal data", err.Error())
	}

	shardIterator := shardIteratorResponse.ShardIterator

	args = gksis.NewArgs()
	args.Add("ShardIterator", shardIterator)
	args.Add("Limit", maxEntries)
	records, err := c.kinesis.GetRecords(args)
	if err != nil {
		return "", tgerrors.NewInternalError("error while reading the internal data", err.Error())
	}

	if err != nil {
		return "", tgerrors.NewInternalError("error while reading the internal data", err.Error())
	}

	if records.NextShardIterator == "" || shardIterator == records.NextShardIterator || len(records.Records) == 0 {
		return nil
	}

	var result []string
	for _, d := range records.Records {
		result = append(result, string(d.GetData()))
	}
	return result, nil
}

func (c *cli) StreamRecords(streamName, consumerName string, maxEntries int) (output chan string, errors chan tgerrors.TGError, done chan struct{}) {

	stream, err := c.DescribeStream(streamName)
	if err != nil {
		errors <- err
		return
	}

	output = make(chan string, len(stream.StreamDescription.Shards*maxEntries))
	errors = make(chan tgerrors.TGError, len(stream.StreamDescription.Shards*maxEntries))
	done = make(chan struct{})
	// Keep track of internal producers and when all of them have quit, we should quit as well
	internalDone := make(chan bool, len(stream.StreamDescription.Shards))

	for idx := range stream.StreamDescription.Shards {
		go func(streamName, shardID string, maxEntries int, output chan<- string, errors chan<- tgerrors.TGError, done chan struct{}) {
			defer func() {
				done <- true
			}()
			args := gksis.NewArgs()
			args.Add("StreamName", streamName)
			args.Add("ShardId", shardID)
			args.Add("ShardIteratorType", "TRIM_HORIZON")
			shardIteratorResponse, err := c.kinesis.GetShardIterator(args)
			if err != nil {
				errors <- tgerrors.NewInternalError("error while reading the internal data", err.Error())
				return
			}

			shardIterator := shardIteratorResponse.ShardIterator

			for {
				args = gksis.NewArgs()
				args.Add("ShardIterator", shardIterator)
				args.Add("Limit", maxEntries)
				records, err := c.kinesis.GetRecords(args)
				if err != nil {
					errors <- tgerrors.NewInternalError("error while reading the internal data", err.Error())
					break
				}

				if err != nil {
					errors <- tgerrors.NewInternalError("error while reading the internal data", err.Error())
					break
				}

				if records.NextShardIterator == "" || shardIterator == records.NextShardIterator {
					errors <- tgerrors.NewInternalError("error while reading the internal data", "shard iterator returned an inconsistent iterator")
					break
				}

				for _, d := range records.Records {
					output <- string(d.GetData())
				}

				shardIterator = records.NextShardIterator
			}
		}(streamName, stream.StreamDescription.Shards[idx].ShardId, maxEntries, output, errors, internalDone)
	}

	go func() {
		i := 0
		for _ := range internalDone {
			i++
			if i == cap(internalDone) {
				break
			}
		}
		close(output)
		close(errors)
		close(done)
	}()
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

func (c *cli) DescribeStream(streamName string) (*gksis.DescribeStreamResp, tgerrors.TGError) {
	type Response struct {
		Descriptor *gksis.DescribeStreamResp
		Error      tgerrors.TGError
	}

	respChan := make(chan Response, 1)

	go func(response chan Response) {
		var err error
		resp := &gksis.DescribeStreamResp{}

		for {
			args := gksis.NewArgs()
			args.Add("StreamName", streamName)
			resp, err = c.kinesis.DescribeStream(args)
			if err != nil {
				response <- Response{
					Descriptor: nil,
					Error:      tgerrors.NewInternalError("failed to read storage", err.Error()),
				}
				return
			}

			if resp.StreamDescription.StreamStatus == "ACTIVE" {
				response <- Response{
					Descriptor: resp,
					Error:      nil,
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}(respChan)

	select {
	case <-time.After(30 * time.Second):
		return nil, tgerrors.NewInternalError("could not connect to the storage in a timely manner", fmt.Sprintf("more than 30 passed in attempting to describe stream %q", streamName))
	case resp := <-respChan:
		return resp.Descriptor, resp.Error
	}
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
