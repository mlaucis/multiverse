/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package logger provides logging functionality to the backend
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	logFormat             = "%s\t%s\t%d\t%s\t%s\t%s\t%+v\t%s\t%s\t%s\t%s\n"
	logResponseTimeFormat = "%s\t%s\t%d\t%s"
	curlGetFormat         = "curl -i %s http://127.0.0.1:8083%s"
	curlPostFormat        = "curl -i -X POST %s -d '%s' http://127.0.0.1:8083%s"
	curlPutFormat         = "curl -i -X PUT %s -d '%s' http://127.0.0.1:8083%s"
	curlDeleteFormat      = "curl -i -X DELETE %s http://127.0.0.1:8083%s"
)

type (
	//LogMsg defines the log message fields
	LogMsg struct {
		RemoteAddr string
		Name       string
		StatusCode int
		Method     string
		RequestURI string
		Headers    http.Header
		Payload    string
		Duration   string
		Message    string
		RawError   string
		Location   string
		Start      time.Time
		End        time.Time `json:"-"`
	}
)

// TGLog is the Tapglue logger
func TGLog(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				if m == nil {
					continue
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
					m.RawError,
				)
			}
		}
	}
}

// TGLogResponseTimes is the Tapglue logger that logs only the names and the response times of requests
func TGLogResponseTimes(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				if m == nil {
					continue
				}

				log.Printf(
					logResponseTimeFormat,
					m.End.Sub(m.Start),
					m.Method,
					m.StatusCode,
					m.Name,
				)
			}
		}
	}
}

// JSONLog provides logging to JSON format
func JSONLog(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				if m == nil {
					continue
				}

				m.Duration = m.End.Sub(m.Start).String()
				m.Headers = getSanitizedHeaders(m.Headers)
				if m.StatusCode < 300 {
					m.Location = ""
				}

				message, err := json.Marshal(m)
				if err != nil {
					log.Print(err)
				} else {
					log.Print(string(message))
				}
			}
		}
	}
}

// TGSilentLog logs all messages to "/dev/null"
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
