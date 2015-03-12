/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package logger provides logging functionality to the backend
package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	logFormat        = "%s\t%s\t%d\t%s\t%s\t%s\t%+v\t%s\t%s\t%s\t%s\n"
	curlGetFormat    = "curl -i %s http://api.tapglue.com%s"
	curlPostFormat   = "curl -i -X POST %s -d '%s' http://api.tapglue.com%s"
	curlPutFormat    = "curl -i -X PUT %s -d '%s' http://api.tapglue.com%s"
	curlDeleteFormat = "curl -i -X DELETE %s http://api.tapglue.com%s"
)

type (
	//LogMsg defines the log message fields
	LogMsg struct {
		RemoteAddr string
		StatusCode int
		Method     string
		RequestURI string
		Name       string
		Message    string
		RawError   error
		Headers    http.Header
		Payload    string
		Location   string
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
				rawError := ""
				if m.RawError != nil {
					rawError = m.RawError.Error()
				}

				log.Printf(
					logFormat,
					m.Message,
					m.RemoteAddr,
					m.StatusCode,
					m.Method,
					m.RequestURI,
					m.Payload,
					getSanitizedHeaders(m.Headers),
					m.Name,
					m.Location,
					m.End.Sub(m.Start),
					rawError,
				)
			}
		}
	}
}

func TGSilentLog(msg chan *LogMsg) {
	for {
		select {
		case _ = <-msg:
			{
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
				case "PUT":
					log.Printf(curlPutFormat, getCurlHeaders(m.Headers), m.Payload, m.RequestURI)
				case "DELETE":
					log.Printf(curlDeleteFormat, getCurlHeaders(m.Headers), m.RequestURI)
				default:
					log.Printf("unexpected curl request\n")
				}
			}
		}
	}
}

// getSanitizedHeaders returns the sanitized request headers
func getSanitizedHeaders(headers http.Header) http.Header {
	// TODO sanitize headers that shouldn't not appear in the logs

	return headers
}
