package main

import (
	"fmt"
	"math/rand"
	"time"

	"flag"
	"os"

	kinesis "github.com/sendgridlabs/go-kinesis"
)

func putRecord(ksis *kinesis.Kinesis, streamName string, partitionID int, producerName string) {
	if !*hasWrite {
		return
	}
	time.Sleep(time.Duration(rand.Intn(100)*100) * time.Millisecond)
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	data := []byte(fmt.Sprintf("Hello AWS Kinesis %s %d. Produced by %s @ %s", streamName, partitionID, producerName, time.Now()))
	partitionKey := fmt.Sprintf("partitionKey-%d", partitionID)
	args.AddRecord(data, partitionKey)
	resp4, err := ksis.PutRecord(args)
	if err != nil {
		fmt.Printf("PutRecord err: %v\n", err)
	} else {
		fmt.Printf("PutRecord: %d %v\n", partitionID, resp4)
	}
}

func getRecords(ksis *kinesis.Kinesis, streamName, ShardID, consumerName, shardSequence string) {
	if !*hasRead {
		return
	}
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", ShardID)
	if shardSequence == "" {
		args.Add("ShardIteratorType", "TRIM_HORIZON")
	} else {
		args.Add("ShardIteratorType", "AT_SEQUENCE_NUMBER")
		args.Add("StartingSequenceNumber", shardSequence)
	}
	resp10, err := ksis.GetShardIterator(args)
	if err != nil {
		panic(err)
	}

	shardIterator := resp10.ShardIterator

	for {
		args = kinesis.NewArgs()
		if consumerName == "consumer1" {
			args.Add("Limit", 1)
		} else if consumerName == "consumer2" {
			args.Add("Limit", 10)
		} else {
			args.Add("Limit", 20)
		}
		args.Add("ShardIterator", shardIterator)
		resp11, err := ksis.GetRecords(args)
		if err != nil {
			fmt.Printf("[%s] There was an error %q\n", consumerName, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if len(resp11.Records) > 0 {
			for _, d := range resp11.Records {
				fmt.Printf("[%s] GetRecords  Data: %v\n", consumerName, string(d.GetData()))
			}
		} else if len(resp11.Records) == 0 {
			fmt.Printf("[%s] Got empty response\n", consumerName)
		} else if resp11.NextShardIterator == "" || shardIterator == resp11.NextShardIterator || err != nil {
			fmt.Printf("[%s] GetRecords ERROR: %v\n", consumerName, err)
			break
		}

		shardIterator = resp11.NextShardIterator
		if consumerName == "consumer1" {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}

func describeStream(ksis *kinesis.Kinesis, streamName string) *kinesis.DescribeStreamResp {
	args := kinesis.NewArgs()
	/*listStreams, err := ksis.ListStreams(args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ListStreams: %v\n", listStreams)*/

	resp := &kinesis.DescribeStreamResp{}
	var err error

	for {
		args = kinesis.NewArgs()
		args.Add("StreamName", streamName)
		resp, err = ksis.DescribeStream(args)
		if err != nil {
			panic(err)
		}
		fmt.Printf("DescribeStream: %v\n", resp)

		if resp.StreamDescription.StreamStatus != "ACTIVE" {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	return resp
}

func setUp(ksis *kinesis.Kinesis, streamName string) {
	err1 := ksis.DeleteStream(streamName)
	if err1 != nil {
		fmt.Printf("DeleteStream ERROR: %v\n", err1)
	}
	time.Sleep(5 * time.Second)
	err := ksis.CreateStream(streamName, 1)
	if err != nil {
		fmt.Printf("CreateStream ERROR: %v\n", err)
	}
}

var hasRead = flag.Bool("read", false, "read")
var hasWrite = flag.Bool("write", false, "write")

func init() {
	flag.Parse()
	if !*hasRead && !*hasWrite {
		fmt.Printf("You must specify if it has at least one of read write")
		os.Exit(64)
	}
}

func main() {
	fmt.Println("Begin")

	stream1Name := "dev"
	// set env variables AWS_ACCESS_KEY and AWS_SECRET_KEY AWS_REGION_NAME
	auth := &kinesis.Auth{
		AccessKey: "AKIAIJ73VAPGUEDRXISA",
		SecretKey: "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",
	}
	ksis := kinesis.New(auth, kinesis.Region{Name: "eu-central-1"})
	//ksis := kinesis.NewWithEndpoint(&auth, kinesis.Region{Name: "eu-central-1"}, "http://127.0.0.1:4567")

	setUp(ksis, stream1Name)

	stream1 := describeStream(ksis, stream1Name)
	fmt.Printf("Stream %s description %#v\n", stream1Name, stream1.StreamDescription.Shards)
	for idx := range stream1.StreamDescription.Shards {
		sequenceNumber := stream1.StreamDescription.Shards[idx].SequenceNumberRange.StartingSequenceNumber
		go getRecords(ksis, stream1Name, stream1.StreamDescription.Shards[idx].ShardId, "s1c1", sequenceNumber)
		go getRecords(ksis, stream1Name, stream1.StreamDescription.Shards[idx].ShardId, "s1c2", sequenceNumber)
	}

	// Wait for user input
	var (
		inputGuess  string
		newConsumer = make(chan bool, 1)
	)
	fmt.Printf("waiting for input ...\n")
	fmt.Scanf("%s\n", &inputGuess)

	go func() {
		<-time.After(7 * time.Second)
		newConsumer <- true
	}()

	var i int
	for {
		select {
		case <-newConsumer:
			for idx := range stream1.StreamDescription.Shards {
				sequenceNumber := stream1.StreamDescription.Shards[idx].SequenceNumberRange.StartingSequenceNumber
				go getRecords(ksis, stream1Name, stream1.StreamDescription.Shards[idx].ShardId, "s1c3", sequenceNumber)
			}
		case <-time.After(time.Duration(1) * time.Second):
			producers := rand.Intn(50)
			for j := 1; j < producers; j++ {
				go putRecord(ksis, stream1Name, i, fmt.Sprintf("s1p%d", j))
			}
			i++
		}
	}
}
