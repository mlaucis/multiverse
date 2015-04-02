/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package core holds the core functionality of the backend
package core

import (
	"github.com/tapglue/backend/v01/storage"

	"gopkg.in/redis.v2"
)

var (
	storageClient *storage.Client
	storageEngine *redis.Client
)

// Init initializes the core package
func Init(engine *storage.Client) {
	storageClient = engine
	storageEngine = engine.Engine()
}
