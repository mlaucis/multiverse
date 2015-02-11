/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package keys

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tapglue/backend/core"
	. "github.com/tapglue/backend/utils"
)

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(requestScope, requestVersion string, r *http.Request) bool {
	signature := r.Header.Get("x-tapglue-signature")
	if signature == "" {
		return false
	}

	if _, err := Base64Decode(signature); err != nil {
		return false
	}

	payload := PeakBody(r).Bytes()
	if Base64Encode(Sha256String(payload)) != r.Header.Get("x-tapglue-payload-hash") {
		return false
	}

	encodedIds := r.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return false
	}

	ids := strings.SplitN(string(decodedIds), ":", 2)
	if len(ids) != 2 {
		return false
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return false
	}

	applicationID, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return false
	}

	application, err := core.ReadApplication(accountID, applicationID)
	if err != nil {
		return false
	}

	signString := generateSigningString(requestScope, requestVersion, r)

	signingKey := generateSigningKey(application.AuthToken, requestVersion, r)

	return r.Header.Get("x-tapglue-signature") == Base64Encode(Sha256String([]byte(signingKey+signString)))
}
