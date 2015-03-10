/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	logFormat      = "%s\t%s\t%s\t%s\t%s\t%+v\t%s\t%s\n"
	curlGetFormat  = "curl -i %s http://localhost:8082%s"
	curlPostFormat = "curl -i -X POST %s -d '%s' http://localhost:8082%s"
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
		Payload    string
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
					m.Message,
					m.RemoteAddr,
					m.Method,
					m.RequestURI,
					m.Payload,
					getSanitizedHeaders(m.Headers),
					m.Name,
					m.End.Sub(m.Start),
				)
			}
		}
	}
}

// getCurlHeaders returns the headers printed as expected by a curl request
func getCurlHeaders(headers http.Header) string {
	result := ""
	for headerName, headerValues := range headers {
		for _, headerValue := range headerValues {
			result += fmt.Sprintf("-H \"%s: %s\" ", headerName, headerValue)
		}
	}
	return result
}

// TGCurlLog is the Tapglue logger which outputs the message as a curl request
func TGCurlLog(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				switch m.Method {
				case "GET":
					log.Printf(curlGetFormat, getCurlHeaders(m.Headers), m.RequestURI)
				case "POST":
					log.Printf(curlPostFormat, getCurlHeaders(m.Headers), m.Payload, m.RequestURI)
				default:
					log.Printf("unexpected curl request\n")
				}
			}
		}
	}
}
