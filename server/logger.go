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
	logFormat = "%s\t%s\t%s\t%+v\t%s\t%s\t%s\n"
)

type (
	//LogMsg defines the log message fields
	LogMsg struct {
		RemoteAddr string
		Method     string
		RequestURI string
		Name       string
		Message    string
		Headers    http.Header
		Start      time.Time
		End        time.Time
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
					m.RemoteAddr,
					m.Method,
					m.RequestURI,
					getSanitizedHeaders(m.Headers),
					m.Name,
					m.End.Sub(m.Start),
					m.Message,
				)
			}
		}
	}
}
