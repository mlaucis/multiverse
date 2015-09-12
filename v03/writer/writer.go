// Package writer defines the interface that the various writers need to implement
package writer

import "github.com/tapglue/backend/errors"

// Writer interface defines the methods needed for a writer
type Writer interface {
	// Execute will provide the main loop logic for the writer
	ProcessMessages(channelName, msg string) errors.Error
}
