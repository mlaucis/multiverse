/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package writer defines the interface that the various writers need to implement
package writer

import "github.com/tapglue/backend/logger"

type (
	// Writer interface defines the methods needed for a writer
	Writer interface {
		// Execute will provide the main loop logic for the writer
		Execute(env string, mainLogChan, errorLogChan chan *logger.LogMsg)
	}
)
