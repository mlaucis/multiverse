/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package core holds the core functionality of the backend
package core

import (
	"github.com/tapglue/backend/v02/storage"

	"github.com/sendgridlabs/go-kinesis"
	"gopkg.in/redis.v2"
)

var (
	storageClient *storage.Client
	redisEngine   *redis.Client
	kinesisEngine *kinesis.Kinesis
)

// Init initializes the core package
func Init(engine *storage.Client) {
	storageClient = engine
	redisEngine = engine.RedisEngine()
	kinesisEngine = engine.KinesisEngine()
}
