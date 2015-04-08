/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package keys

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/context"
	"github.com/tapglue/backend/v01/core"
)

var (
	emptyStringBase64 = utils.Base64Encode(utils.Sha256String(""))
)

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(ctx *context.Context, numKeyParts int) *tgerrors.TGError {
	signature := ctx.R.Header.Get("x-tapglue-signature")
	if signature == "" {
		return tgerrors.NewBadRequestError("failed to verify request signature (1)", "signature failed on 1")
	}

	if _, err := utils.Base64Decode(signature); err != nil {
		return tgerrors.NewBadRequestError("failed to verify request signature (2)", err.Error())
	}

	if ctx.Body != nil {
		if utils.Base64Encode(utils.Sha256String(string(ctx.Body))) != ctx.R.Header.Get("x-tapglue-payload-hash") {
			tgerrors.NewBadRequestError("failed to verify request signature (3)", "signature failed on 3")
		}
	} else {
		if emptyStringBase64 != ctx.R.Header.Get("x-tapglue-payload-hash") {
			tgerrors.NewBadRequestError("failed to verify request signature (4)", "signature failed on 4")
		}
	}

	encodedIds := ctx.R.Header.Get("x-tapglue-id")
	decodedIds, err := utils.Base64Decode(encodedIds)
	if err != nil {
		return tgerrors.NewBadRequestError("failed to verify request signature (5)", err.Error())
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return tgerrors.NewBadRequestError("failed to verify request signature (6)", "signature failed on 6")
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return tgerrors.NewBadRequestError("failed to verify request signature (7)", err.Error())
	}

	authToken := ""
	if numKeyParts == 1 {
		account, err := core.ReadAccount(accountID)
		if err != nil {
			return tgerrors.NewBadRequestError("failed to verify request signature (8)", err.Error())
		}
		authToken = account.AuthToken
	} else {
		applicationID, err := strconv.ParseInt(ids[1], 10, 64)
		if err != nil {
			return tgerrors.NewBadRequestError("failed to verify request signature (9)", err.Error())
		}

		application, er := core.ReadApplication(accountID, applicationID)
		if er != nil {
			return tgerrors.NewBadRequestError("failed to verify request signature (10)", er.Error())
		}

		authToken = application.AuthToken
	}

	signString := generateSigningString(ctx.Scope, ctx.Version, ctx.R)

	signingKey := generateSigningKey(authToken, ctx.Scope, ctx.Version, ctx.R)

	/* TODO Debug content, don't remove unless you want to redo it later ** /
	fmt.Printf("\nURL %s\n", ctx.R.URL.Path)
	fmt.Printf("\nPayload %s - %s \n", ctx.R.Header.Get("x-tapglue-payload-hash"), utils.Base64Encode(utils.Sha256String(ctx.BodyString)))
	fmt.Printf("\nSession %s\n", ctx.R.Header.Get("x-tapglue-session"))
	fmt.Printf("\nSignature parts %s - %s \n", signingKey, signString)
	fmt.Printf("\nSignature %s - %s \n\n", ctx.R.Header.Get("x-tapglue-signature"), utils.Base64Encode(utils.Sha256String(signingKey+signString)))
	/**/
	expectedSignature := utils.Base64Encode(utils.Sha256String(signingKey + signString))
	if ctx.R.Header.Get("x-tapglue-signature") == expectedSignature {
		return nil
	}

	return tgerrors.NewBadRequestError("failed to verify request signature (11)",
		fmt.Sprintf("expected %s got %s", ctx.R.Header.Get("x-tapglue-signature"), expectedSignature))
}
