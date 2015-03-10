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

	. "github.com/tapglue/backend/utils"
)

// addHeaders adds the additional headers to the request before being signed
func addHeaders(accountID, applicationID int64, r *http.Request) error {
	r.Header.Add("x-tapglue-payload-hash", Base64Encode(Sha256String(PeakBody(r).Bytes())))
	if applicationID == 0 {
		r.Header.Add("x-tapglue-id", Base64Encode(fmt.Sprintf("%d", accountID)))
	} else {
		r.Header.Add("x-tapglue-id", Base64Encode(fmt.Sprintf("%d:%d", accountID, applicationID)))
	}

	return nil
}

// canonicalRequest returns the full canonical request and its headers
func canonicalRequest(r *http.Request) []byte {
	// TODO don't allow this in prod environment
	if r.Host == "" {
		r.Host = "localhost:8082"
	}

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
func generateSigningKey(secretKey, scope, requestVersion string, r *http.Request) string {
	key := fmt.Sprintf(
		"tapglue:%s:%s",
		secretKey,
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
				)+scope,
			))+"api",
		)) + requestVersion,
	))
}

// SignRequest runs the signature algorithm on the request and adds the things it's missing
func SignRequest(secretKey, requestScope, requestVersion string, numKeyParts int, r *http.Request) error {
	rawKey, err := Base64Decode(secretKey)
	if err != nil {
		return err
	}

	keyParts := strings.SplitN(string(rawKey), `:`, numKeyParts)
	if len(keyParts) != numKeyParts {
		return fmt.Errorf("not enough key parts")
	}

	accountID, err := strconv.ParseInt(keyParts[0], 10, 64)
	if err != nil {
		return err
	}

	var applicationID int64
	if numKeyParts == 3 {
		applicationID, err = strconv.ParseInt(keyParts[1], 10, 64)
		if err != nil {
			return err
		}
	}

	// Add extra headers
	err = addHeaders(accountID, applicationID, r)
	if err != nil {
		return err
	}

	// Generate signing string
	signString := generateSigningString(requestScope, requestVersion, r)

	// Generate signing key
	signingKey := generateSigningKey(secretKey, requestScope, requestVersion, r)

	// Sign the request
	r.Header.Add("x-tapglue-signature", Base64Encode(Sha256String([]byte(signingKey+signString))))

	return nil
}
