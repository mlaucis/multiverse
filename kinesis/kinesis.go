package main

import (
	"fmt"
	"time"
	"math/rand"

	kinesis "github.com/sendgridlabs/go-kinesis"
)

func putRecord(ksis *kinesis.Kinesis, streamName string, partitionID int) {
	time.Sleep(time.Duration(rand.Intn(100) * 100) * time.Millisecond)
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	data := []byte(fmt.Sprintf("Hello AWS Kinesis %d", partitionID))
	partitionKey := fmt.Sprintf("partitionKey-%d", partitionID)
	args.AddRecord(data, partitionKey)
	resp4, err := ksis.PutRecord(args)
	if err != nil {
		fmt.Printf("PutRecord err: %v\n", err)
	} else {
		fmt.Printf("PutRecord: %d %v\n", partitionID, resp4)
	}
}

func getRecords(ksis *kinesis.Kinesis, streamName, ShardId, consumerName string) {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", ShardId)
	args.Add("ShardIteratorType", "TRIM_HORIZON")
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
			args.Add("Limit", 2)
		} else {
			args.Add("Limit", 10)
		}
		args.Add("ShardIterator", shardIterator)
		resp11, err := ksis.GetRecords(args)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if len(resp11.Records) > 0 {
			for _, d := range resp11.Records {
				fmt.Printf("[%s] GetRecords  Data: %v\n", consumerName, string(d.GetData()))
			}
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
	listStreams, err := ksis.ListStreams(args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ListStreams: %v\n", listStreams)

	resp := &kinesis.DescribeStreamResp{}

	timeout := make(chan bool, 30)
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
			timeout <- true
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

func main() {
	fmt.Println("Begin")

	streamName := "test"
	// set env variables AWS_ACCESS_KEY and AWS_SECRET_KEY AWS_REGION_NAME
	auth := kinesis.NewAuth()
	ksis := kinesis.NewWithEndpoint(&auth, kinesis.Region{Name: "eu-central-1"}, "http://127.0.0.1:4567")

	setUp(ksis, streamName)

	stream := describeStream(ksis, streamName)

	for idx := range stream.StreamDescription.Shards {
		go getRecords(ksis, streamName, stream.StreamDescription.Shards[idx].ShardId, "consumer1")
		go getRecords(ksis, streamName, stream.StreamDescription.Shards[idx].ShardId, "consumer2")
	}

	// Wait for user input
	var (
		inputGuess string
		newConsumer = make(chan bool, 1)
	)
	fmt.Printf("waiting for input: ")
	fmt.Scanf("%s\n", &inputGuess)

	go func() {
		<- time.After(20 * time.Second)
		newConsumer <- true
	}()

	var i int
	for {
		select {
		case <- newConsumer:
			for idx := range stream.StreamDescription.Shards {
				go getRecords(ksis, streamName, stream.StreamDescription.Shards[idx].ShardId, "consumer3")
			}
		case <- time.After(time.Duration(1) * time.Second):
			go putRecord(ksis, streamName, i)
			i++
		}
	}

	fmt.Println("End")
}
