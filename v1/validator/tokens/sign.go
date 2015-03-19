/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package tokens

import (
	"net/http"
)

// SignRequest runs the signature algorithm on the request and adds the things it's missing
func SignRequest(secretKey, requestScope, requestVersion string, numKeyParts int, r *http.Request) error {
	r.Header.Add("x-tapglue-app-key", secretKey)

	return nil
}
