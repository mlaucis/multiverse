package kinesis_test

import (
	"fmt"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/kinesis"
)

func main() {
	svc := kinesis.New(&aws.Config{
		Credentials: aws.Creds("BAD_KEY", "BAD_SECRETE_KEY", "SESSION_TOKEN"),
		Endpoint: "http://127.0.0.1:4567",
		Region: "eu-central-1",
	})

	params := &kinesis.CreateStreamInput{
		ShardCount: aws.Long(1),              // Required
		StreamName: aws.String("StreamName"), // Required
	}
	resp, err := svc.CreateStream(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(resp))

	time.Sleep(time.Duration(150) * time.Millisecond)

	paramsList := &kinesis.ListStreamsInput{
		ExclusiveStartStreamName: aws.String("StreamName"),
		Limit: aws.Long(1),
	}
	respList, err := svc.ListStreams(paramsList)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(respList))

	time.Sleep(time.Duration(150) * time.Millisecond)

	paramsDescribe := &kinesis.DescribeStreamInput{
		StreamName:            aws.String("StreamName"), // Required
		ExclusiveStartShardID: aws.String("ShardId"),
		Limit: aws.Long(1),
	}
	respDescribe, err := svc.DescribeStream(paramsDescribe)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(respDescribe))


	time.Sleep(time.Duration(150) * time.Millisecond)

	paramsDelete := &kinesis.DeleteStreamInput{
		StreamName: aws.String("StreamName"), // Required
	}
	respDelete, err := svc.DeleteStream(paramsDelete)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(respDelete))
}
