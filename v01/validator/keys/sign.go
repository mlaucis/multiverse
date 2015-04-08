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

	"github.com/tapglue/backend/utils"
)

// addHeaders adds the additional headers to the request before being signed
func addHeaders(accountID, applicationID int64, r *http.Request) error {
	r.Header.Add("x-tapglue-payload-hash", utils.Base64Encode(utils.Sha256String(utils.PeakBody(r).String())))
	if applicationID == 0 {
		r.Header.Add("x-tapglue-id", utils.Base64Encode(fmt.Sprintf("%d", accountID)))
	} else {
		r.Header.Add("x-tapglue-id", utils.Base64Encode(fmt.Sprintf("%d:%d", accountID, applicationID)))
	}

	return nil
}

// canonicalRequest returns the full canonical request and its headers
func canonicalRequest(r *http.Request) string {
	result := fmt.Sprintf(
		"%s\n%s\nhost:%s\nx-tapglue-date:%s\nx-tapglue-payload-hash:%s\nx-tapglue-id:%s",
		r.Method,
		r.URL.Path,
		r.Host,
		r.Header.Get("x-tapglue-date"),
		r.Header.Get("x-tapglue-payload-hash"),
		r.Header.Get("x-tapglue-id"),
	)

	return result
}

// generateSigningString returns the string used to
func generateSigningString(scope, requestVersion string, r *http.Request) string {
	return requestVersion + "\n" +
		r.Header.Get("x-tapglue-date") + "\n" +
		getScope(r.Header.Get("x-tapglue-date"), scope, requestVersion) + "\n" +
		utils.Base64Encode(utils.Sha256String(canonicalRequest(r)))
}

// generateSigningKey returns the key used to sign the request
func generateSigningKey(secretKey, scope, requestVersion string, r *http.Request) string {
	key := fmt.Sprintf(
		"tapglue:%s:%s",
		secretKey,
		r.Header.Get("x-tapglue-date"),
	)

	//log.Printf("sign:\tkey\t%s", key)
	key = utils.Base64Encode(utils.Sha256String(key))
	//log.Printf("sign:\tsha key\t%s", key)
	key = utils.Base64Encode(utils.Sha256String(key + r.Header.Get("x-tapglue-session")))
	//log.Printf("sign:\tkey + session\t%s", key)
	key = utils.Base64Encode(utils.Sha256String(key + scope))
	//log.Printf("sign:\tkey + scope\t%s", key)
	key = utils.Base64Encode(utils.Sha256String(key + "api"))
	//log.Printf("sign:\tkey + \"api\"\t%s", key)

	key = utils.Base64Encode(utils.Sha256String(key + requestVersion))
	//log.Printf("sign:\tkey + requestVersion\t%s", key)

	return key
}

// SignRequest runs the signature algorithm on the request and adds the things it's missing
func SignRequest(secretKey, requestScope, requestVersion string, numKeyParts int, r *http.Request) error {
	rawKey, err := utils.Base64Decode(secretKey)
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
	r.Header.Add("x-tapglue-signature", utils.Base64Encode(utils.Sha256String(signingKey+signString)))

	return nil
}
