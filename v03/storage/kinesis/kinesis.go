// Package kinesis provides the AWS Kinesis needed functions for kinesis
package kinesis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"

	gksis "github.com/sendgridlabs/go-kinesis"
)

type (
	// Client defines the interface for Kinesis
	Client interface {
		// SetupStreams creates the needed Kinesis streams
		SetupStreams([]string) error

		// PutRecord sends a new record to a Kinesis stream
		PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, errors.Error)

		// PackAndPutRecords will pack the record to minimize the number of needed streams
		//
		// Use UnpackRecord to get the original target stream and message
		PackAndPutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, errors.Error)

		// GetRecords returns at most the specified number of records from the desired stream / shard
		GetRecords(streamName, shardID, consumerName string, maxEntries int) ([]string, errors.Error)

		// StreamRecords will stream all the records it can from from all the shards of a stream
		// To stop it, just close the output channel
		StreamRecords(streamName, consumerName, consumerPosition string, maxEntries int) (<-chan string, <-chan string, chan errors.Error, <-chan struct{})

		// UnpackRecord takes a record and unpacks it then returns the stream name and the original message as a string
		UnpackRecord(message string) (streamName, unpackedMessage string, err errors.Error)

		// DescribeStream will return the Kinesis stream descriptiont that AWS has if the stream is active.
		// If the stream is not active then it will return an error
		// It will also timeout after 30 seconds (one try per second)
		DescribeStream(streamName string) (*gksis.DescribeStreamResp, errors.Error)

		// TeardownStreams destroys the streams from Kinesis
		TeardownStreams(streamsName []string) error

		// Datastore returns the Kinesis client
		Datastore() *gksis.Kinesis
	}

	cli struct {
		kinesis          *gksis.Kinesis
		packedStreamName string
	}

	packedPayload struct {
		StreamName string `json:"stream_name"`
		Message    string `json:"message"`
	}
)

// These are the names of the Kinesis streams that we can use across various platforms
const (
	StreamAccountUpdate           = "v03_account_update"
	StreamAccountDelete           = "v03_account_delete"
	StreamAccountUserCreate       = "v03_account_user_create"
	StreamAccountUserUpdate       = "v03_account_user_update"
	StreamAccountUserDelete       = "v03_account_user_delete"
	StreamApplicationCreate       = "v03_application_create"
	StreamApplicationUpdate       = "v03_application_update"
	StreamApplicationDelete       = "v03_application_delete"
	StreamApplicationUserUpdate   = "v03_application_user_update"
	StreamApplicationUserDelete   = "v03_application_user_delete"
	StreamConnectionCreate        = "v03_connection_create"
	StreamConnectionUpdate        = "v03_connection_update"
	StreamConnectionDelete        = "v03_connection_delete"
	StreamConnectionConfirm       = "v03_connection_confirm"
	StreamConnectionSocialConnect = "v03_connection_social_connect"
	StreamConnectionAutoConnect   = "v03_connection_auto_connect"
	StreamEventCreate             = "v03_event_create"
	StreamEventUpdate             = "v03_event_update"
	StreamEventDelete             = "v03_event_delete"
)

var (
	// Streams defines the array of all streams currently defined by the application
	Streams = []string{
		StreamAccountUpdate,
		StreamAccountDelete,
		StreamAccountUserCreate,
		StreamAccountUserUpdate,
		StreamAccountUserDelete,
		StreamApplicationCreate,
		StreamApplicationUpdate,
		StreamApplicationDelete,
		StreamApplicationUserUpdate,
		StreamApplicationUserDelete,
		StreamConnectionCreate,
		StreamConnectionUpdate,
		StreamConnectionDelete,
		StreamConnectionConfirm,
		StreamConnectionSocialConnect,
		StreamConnectionAutoConnect,
		StreamEventCreate,
		StreamEventUpdate,
		StreamEventDelete,
	}
)

func (c *cli) PutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, errors.Error) {
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.AddRecord(payload, partitionKey)
	resp, err := c.kinesis.PutRecord(args)
	if err != nil {
		return nil, errors.NewInternalError(0, "failed to execute operation (1)", err.Error())
	}

	return resp, nil
}

func (c *cli) PackAndPutRecord(streamName, partitionKey string, payload []byte) (*gksis.PutRecordResp, errors.Error) {
	packedPayload := packedPayload{
		StreamName: streamName,
		Message:    utils.Base64Encode(string(payload)),
	}
	myPayload, err := json.Marshal(packedPayload)
	if err != nil {
		return nil, errors.NewInternalError(0, "failed to generate the message", err.Error())
	}
	return c.PutRecord(c.packedStreamName, partitionKey, myPayload)
}

func (c *cli) GetRecords(streamName, shardID, consumerName string, maxEntries int) ([]string, errors.Error) {
	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", shardID)
	// TODO this one should actually be retrieved from REDIS or other place where we can store the current iterator
	// so that we can have resume support for it
	args.Add("ShardIteratorType", "LATEST")
	shardIteratorResponse, err := c.kinesis.GetShardIterator(args)
	if err != nil {
		return []string{}, errors.NewInternalError(0, "error while reading the internal data", err.Error())
	}

	shardIterator := shardIteratorResponse.ShardIterator

	args = gksis.NewArgs()
	args.Add("ShardIterator", shardIterator)
	args.Add("Limit", maxEntries)
	records, err := c.kinesis.GetRecords(args)
	if err != nil {
		return []string{}, errors.NewInternalError(0, "error while reading the internal data", err.Error())
	}

	if err != nil {
		return []string{}, errors.NewInternalError(0, "error while reading the internal data", err.Error())
	}

	if records.NextShardIterator == "" || shardIterator == records.NextShardIterator {
		return []string{}, errors.NewInternalError(0, "error while reading the internal data", "malformed pointer received")
	}

	if len(records.Records) == 0 {
		return []string{}, nil
	}

	var result []string
	for _, d := range records.Records {
		result = append(result, string(d.GetData()))
	}
	return result, nil
}

func (c *cli) StreamRecords(streamName, consumerName, consumerPosition string, maxEntries int) (<-chan string, <-chan string, chan errors.Error, <-chan struct{}) {

	output := make(chan string, 10*maxEntries)
	sequenceNumber := make(chan string, 10*maxEntries)
	errs := make(chan errors.Error, 10)
	done := make(chan struct{})

	stream, err := c.DescribeStream(streamName)
	if err != nil {
		errs <- err
		close(done)
		return output, sequenceNumber, errs, done
	}

	// Keep track of internal producers and when all of them have quit, we should quit as well
	internalDone := make(chan bool, len(stream.StreamDescription.Shards))

	for idx := range stream.StreamDescription.Shards {
		go c.streamShard(consumerPosition, streamName, stream.StreamDescription.Shards[idx].ShardId, maxEntries, output, sequenceNumber, errs, internalDone)
	}

	go func() {
		i := 0
		for _ = range internalDone {
			i++
			if i == cap(internalDone) {
				break
			}
		}
		close(output)
		close(sequenceNumber)
		close(errs)
		close(done)
	}()

	return output, sequenceNumber, errs, done
}

func (c *cli) streamShard(consumerPosition, streamName, shardID string, maxEntries int, output, sequenceNumber chan<- string, errs chan<- errors.Error, done chan bool) {
	defer func() {
		done <- true
	}()

	args := gksis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", shardID)
	if consumerPosition == "" {
		args.Add("ShardIteratorType", "LATEST")
	} else {
		args.Add("ShardIteratorType", "AFTER_SEQUENCE_NUMBER")
		args.Add("StartingSequenceNumber", consumerPosition)
	}
	shardIteratorResponse, err := c.kinesis.GetShardIterator(args)
	if err != nil {
		errs <- errors.NewInternalError(0, "error while reading the internal data", err.Error())
		return
	}

	shardIterator := shardIteratorResponse.ShardIterator

	for {
		args = gksis.NewArgs()
		args.Add("ShardIterator", shardIterator)
		args.Add("Limit", maxEntries)
		records, err := c.kinesis.GetRecords(args)
		if err != nil {
			errs <- errors.NewInternalError(0, "error while reading the internal data", err.Error())
			continue
		}

		if records.NextShardIterator == "" || shardIterator == records.NextShardIterator {
			errs <- errors.NewInternalError(0, "error while reading the internal data", "shard iterator returned an inconsistent iterator")
			continue
		}

		for _, d := range records.Records {
			output <- string(d.GetData())
			sequenceNumber <- d.SequenceNumber
		}

		shardIterator = records.NextShardIterator
		time.Sleep(1 * time.Second)
	}
}

func (c *cli) UnpackRecord(message string) (streamName, unpackedMessage string, err errors.Error) {
	unpackedPayload := packedPayload{}
	er := json.Unmarshal([]byte(message), &unpackedPayload)
	if er != nil {
		return "", "", errors.NewInternalError(0, "failed to receive the message", er.Error())
	}
	unpackedMessage, er = utils.Base64Decode(unpackedPayload.Message)
	if er != nil {
		return "", "", errors.NewInternalError(0, "failed to decode the received message", er.Error())
	}
	return unpackedPayload.StreamName, unpackedMessage, nil
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

func (c *cli) DescribeStream(streamName string) (*gksis.DescribeStreamResp, errors.Error) {
	type Response struct {
		Descriptor *gksis.DescribeStreamResp
		Error      errors.Error
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
					Error:      errors.NewInternalError(0, "failed to read storage", err.Error()),
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
		return nil, errors.NewInternalError(0, "could not connect to the storage in a timely manner", fmt.Sprintf("more than 30 passed in attempting to describe stream %q", streamName))
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
func New(authKey, secretKey, region, env, packedStreamName string) Client {
	auth := gksis.NewAuth(authKey, secretKey)

	return &cli{
		kinesis:          gksis.New(auth, region),
		packedStreamName: packedStreamName,
	}
}

// NewWithEndpoint returns a new testing-enabled client
func NewWithEndpoint(authKey, secretKey, region, endpoint, env, packedStreamName string) Client {
	auth := gksis.NewAuth(authKey, secretKey)

	return &cli{
		kinesis:          gksis.NewWithEndpoint(auth, region, endpoint),
		packedStreamName: packedStreamName,
	}
}
