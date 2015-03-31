/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package tokens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	. "github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/core"
)

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(ctx *context.Context, numKeyParts int) *tgerrors.TGError {
	encodedIds := ctx.R.Header.Get("x-tapglue-app-key")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return tgerrors.NewBadRequestError("application key validation failed (1)", err.Error())
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return tgerrors.NewBadRequestError("application key validation failed (2)", fmt.Sprintf("app key parts count mismatch expected %d got %d", numKeyParts, len(ids)))
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return tgerrors.NewBadRequestError("application key validation failed (3)", "failed to parse account id")
	}

	if accountID != ctx.AccountID {
		return tgerrors.NewBadRequestError("application key validation failed (4)", fmt.Sprintf("account id mismatch expected %d got %d", ctx.AccountID, accountID))
	}

	authToken := ""
	if numKeyParts == 1 {
		account, er := core.ReadAccount(accountID)
		if er != nil {
			return tgerrors.NewBadRequestError("application key validation failed (5)", er.Error())
		}
		authToken = account.AuthToken
	} else {
		applicationID, er := strconv.ParseInt(ids[1], 10, 64)
		if er != nil {
			return tgerrors.NewBadRequestError("application key validation failed (6)", er.Error())
		}

		if applicationID != ctx.ApplicationID {
			return tgerrors.NewBadRequestError("application key validation failed (7)", fmt.Sprintf("app id mismatch expected %d got %d", ctx.ApplicationID, applicationID))
		}

		application, err := core.ReadApplication(accountID, applicationID)
		if err != nil {
			return tgerrors.NewBadRequestError("application key validation failed (8)", er.Error())
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
		return nil
	}

	return tgerrors.NewBadRequestError("application key validation failed (9)\nsignature mismatch",
		fmt.Sprintf("app key mismatch expected %s got %s", ctx.R.Header.Get("x-tapglue-app-key"), authToken))
}
