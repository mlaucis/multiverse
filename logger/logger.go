// Package logger provides logging functionality to the backend
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"regexp"
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
		RemoteAddr string      `json:"remote_addr,omitempty"`
		Name       string      `json:"name,omitempty"`
		StatusCode int         `json:"status_code,omitempty"`
		Method     string      `json:"method,omitempty"`
		RequestURI string      `json:"request_uri,omitempty"`
		Headers    http.Header `json:"headers,omitempty"`
		Payload    string      `json:"payload,omitempty"`
		Duration   string      `json:"duration,omitempty"`
		Message    string      `json:"message,omitempty"`
		RawError   string      `json:"raw_error,omitempty"`
		Location   string      `json:"location,omitempty"`
		Start      time.Time   `json:"start,omitempty"`
		End        time.Time   `json:"-"`
	}
)

var (
	// This won't catch all the passwords, think passwords that have double-quotes in them
	// but we shouldn't never have those here anyway, clients should never send us passwords
	// in clear, right? Right? RIGHT?
	passwordRE = regexp.MustCompile(`("password": ".*?")`)
	passwordREReplacement = []byte(`"password": ""`)
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

				m.Payload = string(passwordRE.ReplaceAll([]byte(m.Payload), passwordREReplacement))

				log.Printf(
					logFormat,
					m.Message,
					m.RemoteAddr,
					m.StatusCode,
					m.Method,
					m.RequestURI,
					m.Payload,
					m.Headers,
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
				if m.StatusCode < 300 {
					m.Location = ""
				}
				m.Payload = string(passwordRE.ReplaceAll([]byte(m.Payload), passwordREReplacement))

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
