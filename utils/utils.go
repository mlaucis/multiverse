/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// We have to have our own type so that we can break what go forces us to do
type noCloseReaderCloser struct {
	*bytes.Buffer
}

// We should do some closing here but then again, that's what we want to prevent
func (m noCloseReaderCloser) Close() error {
	return nil
}

// PeakBody allows us to look at the request body and get the values without closing the body
func PeakBody(r *http.Request) *bytes.Buffer {
	buf, _ := ioutil.ReadAll(r.Body)
	buff := noCloseReaderCloser{bytes.NewBuffer(buf)}
	r.Body = noCloseReaderCloser{bytes.NewBuffer(buf)}
	return buff.Buffer
}

// sha256String takes a byte slice and returns its sha256 checksum value
func Sha256String(value []byte) string {
	hasher := sha256.New()
	hasher.Write(value)

	return string(hasher.Sum(nil))
}

// base64Encode encodes a string in base64
func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// base64Decode decodes a string from base64 to the decoded version
func Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
