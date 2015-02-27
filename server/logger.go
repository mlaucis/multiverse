/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"log"
	"net/http"
	"time"
)

const (
	logFormat = "%s\t%s\t%+v\t%s\t%s\t%s\n"
)

type (
	//LogMsg defines the log message fields
	LogMsg struct {
		method     string
		requestURI string
		name       string
		message    string
		headers    http.Header
		start      time.Time
		end        time.Time
	}
)

// TGLog is the Tapglue logger
func TGLog(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				log.Printf(
					logFormat,
					m.method,
					m.requestURI,
					getSanitizedHeaders(m.headers),
					m.name,
					m.end.Sub(m.start),
					m.message,
				)
			}
		}
	}
}
