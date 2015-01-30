/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/storage"

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
