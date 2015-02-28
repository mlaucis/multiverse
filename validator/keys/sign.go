/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package keys provides a way to sign and check a request
package keys

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/tapglue/backend/utils"
)

// addHeaders adds the additional headers to the request before being signed
func addHeaders(accountID, applicationID int64, r *http.Request) error {
	r.Header.Add("x-tapglue-date", time.Now().Format(time.RFC3339))
	r.Header.Add("x-tapglue-payload-hash", Base64Encode(Sha256String(PeakBody(r).Bytes())))
	r.Header.Add("x-tapglue-id", Base64Encode(fmt.Sprintf("%d:%d", accountID, applicationID)))

	return nil
}

// canonicalRequest returns the full canonical request and its headers
func canonicalRequest(r *http.Request) []byte {
	req := fmt.Sprintf(
		"%s\n%s\nhost:%s\nx-tapglue-date:%s\nx-tapglue-payload-hash:%s\nx-tapglue-id:%s",
		r.Method,
		r.URL.Path,
		r.Host,
		r.Header.Get("x-tapglue-date"),
		r.Header.Get("x-tapglue-payload-hash"),
		r.Header.Get("x-tapglue-id"),
	)

	return []byte(req)
}

// generateSigningString returns the string used to
func generateSigningString(scope, requestVersion string, r *http.Request) string {
	return requestVersion + "\n" +
		r.Header.Get("x-tapglue-date") + "\n" +
		getScope(r.Header.Get("x-tapglue-date"), scope, requestVersion) + "\n" +
		Base64Encode(Sha256String(canonicalRequest(r)))
}

// generateSigningKey returns the key used to sign the request
func generateSigningKey(applicationSecretKey, requestVersion string, r *http.Request) string {
	key := fmt.Sprintf(
		"tapglue:%s:%s",
		applicationSecretKey,
		r.Header.Get("x-tapglue-date"),
	)
	return Sha256String([]byte(
		Sha256String([]byte(
			Sha256String([]byte(
				Sha256String(
					[]byte(
						Sha256String([]byte(key))+
						r.Header.Get("x-tapglue-session"),
					),
				)+"user/log",
			))+"api",
		)) + requestVersion,
	))
}

// SignRequest runs the signature algorithm on the request and adds the things it's missing
func SignRequest(applicationSecretKey, requestScope, requestVersion string, r *http.Request) error {
	rawKey, err := Base64Decode(applicationSecretKey)
	if err != nil {
		return err
	}

	keyParts := strings.SplitN(string(rawKey), `:`, 3)
	if len(keyParts) != 3 {
		return fmt.Errorf("not enough key parts")
	}

	accountID, err := strconv.ParseInt(keyParts[0], 10, 64)
	if err != nil {
		return err
	}

	applicationID, err := strconv.ParseInt(keyParts[1], 10, 64)
	if err != nil {
		return err
	}

	// Add extra headers
	err = addHeaders(accountID, applicationID, r)
	if err != nil {
		return err
	}

	// Generate signing string
	signString := generateSigningString(requestScope, requestVersion, r)

	// Generate signing key
	signingKey := generateSigningKey(applicationSecretKey, requestVersion, r)

	// Sign the request
	r.Header.Add("x-tapglue-signature", Base64Encode(Sha256String([]byte(signingKey+signString))))

	return nil
}
