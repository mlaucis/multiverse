package main

import (
	"fmt"
	"time"
	"math/rand"

	kinesis "github.com/sendgridlabs/go-kinesis"
)

func putRecord(ksis *kinesis.Kinesis, streamName string, i int) {
	time.Sleep(time.Duration(rand.Intn(100) * 100) * time.Millisecond)
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	data := []byte(fmt.Sprintf("Hello AWS Kinesis %d", i))
	partitionKey := fmt.Sprintf("partitionKey-%d", i)
	args.AddRecord(data, partitionKey)
	resp4, err := ksis.PutRecord(args)
	if err != nil {
		fmt.Printf("PutRecord err: %v\n", err)
	} else {
		fmt.Printf("PutRecord: %d %v\n", i, resp4)
	}
}

func getRecords(ksis *kinesis.Kinesis, streamName, ShardId, consumerName string) {
	args := kinesis.NewArgs()
	args.Add("StreamName", streamName)
	args.Add("ShardId", ShardId)
	args.Add("ShardIteratorType", "TRIM_HORIZON")
	resp10, _ := ksis.GetShardIterator(args)

	shardIterator := resp10.ShardIterator

	for {
		args = kinesis.NewArgs()
		args.Add("ShardIterator", shardIterator)
		resp11, err := ksis.GetRecords(args)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		if len(resp11.Records) > 0 {
			fmt.Printf("[%s] GetRecords Data BEGIN\n", consumerName)
			for _, d := range resp11.Records {
				fmt.Printf("[%s] GetRecords  Data: %v\n", consumerName, string(d.GetData()))
			}
			fmt.Printf("[%s] GetRecords Data END\n", consumerName)
		} else if resp11.NextShardIterator == "" || shardIterator == resp11.NextShardIterator || err != nil {
			fmt.Printf("[%s] GetRecords ERROR: %v\n", consumerName, err)
			break
		}

		shardIterator = resp11.NextShardIterator
		if consumerName == "consumer1" {
			time.Sleep(1000 * time.Millisecond)
		} else {
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func describeStream(ksis *kinesis.Kinesis, streamName string) *kinesis.DescribeStreamResp {
	args := kinesis.NewArgs()
	resp2, _ := ksis.ListStreams(args)
	fmt.Printf("ListStreams: %v\n", resp2)

	resp := &kinesis.DescribeStreamResp{}

	timeout := make(chan bool, 30)
	for {
		args = kinesis.NewArgs()
		args.Add("StreamName", streamName)
		resp, _ = ksis.DescribeStream(args)
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

	for _, shard := range stream.StreamDescription.Shards {
		go getRecords(ksis, streamName, shard.ShardId, "consumer1")
		go getRecords(ksis, streamName, shard.ShardId, "consumer2")
	}

	// Wait for user input
	var inputGuess string
	fmt.Printf("waiting for input: ")
	fmt.Scanf("%s\n", &inputGuess)

	var i int
	for {
		select {
		case <- time.After(time.Duration(1) * time.Second):
			go putRecord(ksis, streamName, i)
			i++
		}
	}

	fmt.Println("End")
}
