/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package kinesis provides the AWS Kinesis needed functions for kinesis
package kinesis

import gksis "github.com/sendgridlabs/go-kinesis"

type (
	cli struct {
		client *gksis.Kinesis
	}
)

var (
	kinesisClient *cli
)

// Init initializes the redis client
func Init(authKey, secretKey, region string) {
	auth := &gksis.Auth{
		AccessKey: authKey,
		SecretKey: secretKey,
	}
	kinesisClient = &cli{
		client: gksis.New(auth, gksis.Region{Name: region}),
	}
}

// Client returns the redis client
func Client() *gksis.Kinesis {
	return kinesisClient.client
}
