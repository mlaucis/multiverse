/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package tokens

import (
	"strconv"
	"strings"

	"fmt"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	. "github.com/tapglue/backend/utils"
)

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(ctx *context.Context, numKeyParts int) error {
	encodedIds := ctx.R.Header.Get("x-tapglue-app-key")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return fmt.Errorf("signature failed on 2")
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return fmt.Errorf("signature failed on 2")
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return fmt.Errorf("signature failed on 3")
	}

	if accountID != ctx.AccountID {
		return fmt.Errorf("account ids mismatch")
	}

	authToken := ""
	if numKeyParts == 1 {
		account, err := core.ReadAccount(accountID)
		if err != nil {
			return fmt.Errorf("signature failed on 4")
		}
		authToken = account.AuthToken
	} else {
		applicationID, err := strconv.ParseInt(ids[1], 10, 64)
		if err != nil {
			return fmt.Errorf("signature failed on 5")
		}

		if applicationID != ctx.ApplicationID {
			return fmt.Errorf("application ids mismatch")
		}

		application, err := core.ReadApplication(accountID, applicationID)
		if err != nil {
			return fmt.Errorf("signature failed on 6")
		}

		authToken = application.AuthToken
	}

	/* TODO Debug content, don't remove unless you want to redo it later
	   fmt.Printf("\nPayload %s - %s \n", r.Header.Get("x-tapglue-payload-hash"), Base64Encode(Sha256String(payload)))
	   fmt.Printf("\nSession %s\n", r.Header.Get("x-tapglue-session"))
	   fmt.Printf("\nSignature parts %s - %s \n", Base64Encode(signingKey), Base64Encode(signString))
	   fmt.Printf("\nSignature %s - %s \n\n", r.Header.Get("x-tapglue-signature"), Base64Encode(Sha256String([]byte(signingKey+signString))))
	*/

	if ctx.R.Header.Get("x-tapglue-app-key") != authToken {
		return fmt.Errorf("signature failed on 7")
	}
	return nil
}
