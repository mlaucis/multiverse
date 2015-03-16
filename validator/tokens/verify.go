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
func VerifyRequest(ctx *context.Context, numKeyParts int) (string, error) {
	encodedIds := ctx.R.Header.Get("x-tapglue-app-key")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "application key validation failed (1)", err
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return "application key validation failed (2)", fmt.Errorf("app key parts count mismatch expected %d got %d", numKeyParts, len(ids))
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return "application key validation failed (3)", fmt.Errorf("failed to parse account id")
	}

	if accountID != ctx.AccountID {
		return "application key validation failed (4)", fmt.Errorf("account id mismatch expected %d got %d", ctx.AccountID, accountID)
	}

	authToken := ""
	if numKeyParts == 1 {
		account, err := core.ReadAccount(accountID)
		if err != nil {
			return "application key validation failed (5)", err
		}
		authToken = account.AuthToken
	} else {
		applicationID, err := strconv.ParseInt(ids[1], 10, 64)
		if err != nil {
			return "application key validation failed (6)", err
		}

		if applicationID != ctx.ApplicationID {
			return "application key validation failed (7)", fmt.Errorf("app id mismatch expected %d got %d", ctx.ApplicationID, applicationID)
		}

		application, err := core.ReadApplication(accountID, applicationID)
		if err != nil {
			return "application key validation failed (8)", err
		}

		authToken = application.AuthToken
	}

	/* TODO Debug content, don't remove unless you want to redo it later
	   fmt.Printf("\nPayload %s - %s \n", r.Header.Get("x-tapglue-payload-hash"), Base64Encode(Sha256String(payload)))
	   fmt.Printf("\nSession %s\n", r.Header.Get("x-tapglue-session"))
	   fmt.Printf("\nSignature parts %s - %s \n", Base64Encode(signingKey), Base64Encode(signString))
	   fmt.Printf("\nSignature %s - %s \n\n", r.Header.Get("x-tapglue-signature"), Base64Encode(Sha256String([]byte(signingKey+signString))))
	*/

	if ctx.R.Header.Get("x-tapglue-app-key") == authToken {
		return "", nil
	}

	return "application key validation failed (9)\nsignature mismatch", fmt.Errorf("app key mismatch expected %s got %s", ctx.R.Header.Get("x-tapglue-app-key"), authToken)
}
