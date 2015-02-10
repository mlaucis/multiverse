package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/kr/pretty"
)

var (
	accountId               int64 = 1
	applicationID           int64 = 1
	applicationOriginalName       = "Demo app"
	applicationCreatedAt          = time.Now().Format(time.RFC3339)
	applicationSecretKey          = "application_secret_key"
	requestVersion                = "tg_0.1_request"
)

func sha256String(value []byte) string {
	hasher := sha256.New()
	hasher.Write(value)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func generateSecretKey() string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf(
		"%d%d%s%s",
		accountId,
		applicationID,
		applicationOriginalName,
		applicationCreatedAt,
	)))

	return base64.URLEncoding.EncodeToString([]byte(
		fmt.Sprintf(
			"%d:%d:%s",
			accountId,
			applicationID,
			string(hasher.Sum(nil)),
		)))
}

func addHeaders(r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	r.Header.Add("x-tagplue-date", time.Now().Format(time.RFC3339))
	r.Header.Add("x-tagplue-payload-hash", sha256String([]byte(body)))
	r.Header.Add("x-tapglue-id", fmt.Sprintf("%d", applicationID))
}

func getScope(date, scope string) string {
	return date + "/" + scope + "/" + requestVersion
}

func canonicalRequest(r *http.Request) []byte {
	req := fmt.Sprintf(
		"%s\n%s\nhost:%s\nx-tagplue-date:%s\nx-tagplue-payload-hash:%s\nx-tapglue-id:%s",
		r.Method,
		r.URL.Path,
		r.Host,
		r.Header.Get("x-tapglue-date"),
		r.Header.Get("x-tapglue-payload-hash"),
		r.Header.Get("x-tapglue-id"),
	)

	return []byte(req)
}

func generateSigningString(r *http.Request, scope string) string {
	return requestVersion + "\n" +
		r.Header.Get("x-tapglue-date") + "\n" +
		getScope(r.Header.Get("x-tapglue-date"), scope) + "\n" +
		sha256String(canonicalRequest(r))
}

func generateSigningKey(r *http.Request) string {
	return sha256String([]byte(
		sha256String([]byte(
			sha256String([]byte(
				sha256String([]byte(
					"tapglue:"+applicationSecretKey+":"+r.Header.Get("x-tapglue-date"),
				))+"user/log",
			))+"api",
		)) + requestVersion,
	))
}

func signRequest(r *http.Request, scope string) {
	// Add extra headers
	addHeaders(r)

	// Generate signing string
	signString := generateSigningString(r, scope)

	// Generate signing key
	signningKey := generateSigningKey(r)

	// Sign the request
	r.Header.Add("x-tapglue-signature", sha256String([]byte(signningKey+signString)))
}

func main() {

	applicationSecretKey = generateSecretKey()

	jsonStr := []byte(`{"username": "florin", "password": "passwd"}`)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.tagpglue.com/acc/%d/app/%d/user/login", accountId, applicationID),
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		panic(err)
	}
	signRequest(req, "user/login")

	pretty.Printf("%# v\n", req)

	decodedKey, _ := base64.URLEncoding.DecodeString(applicationSecretKey)
	fmt.Printf("\nAccount key: %#v\nDecoded key: %#v", applicationSecretKey, string(decodedKey))
}
