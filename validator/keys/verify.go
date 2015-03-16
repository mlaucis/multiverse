/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package keys

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	. "github.com/tapglue/backend/utils"
)

var (
	emptyStringBase64 = Base64Encode(Sha256String(""))
)

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(ctx *context.Context, numKeyParts int) (string, error) {
	signature := ctx.R.Header.Get("x-tapglue-signature")
	if signature == "" {
		return "failed to verify request signature (1)", fmt.Errorf("signature failed on 1")
	}

	if _, err := Base64Decode(signature); err != nil {
		return "failed to verify request signature (2)", err
	}

	if ctx.Body != nil {
		if Base64Encode(Sha256String(ctx.Body.String())) != ctx.R.Header.Get("x-tapglue-payload-hash") {
			return "failed to verify request signature (3)", fmt.Errorf("signature failed on 3")
		}
	} else {
		if emptyStringBase64 != ctx.R.Header.Get("x-tapglue-payload-hash") {
			return "failed to verify request signature (4)", fmt.Errorf("signature failed on 3")
		}
	}

	encodedIds := ctx.R.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "failed to verify request signature (5)", err
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return "failed to verify request signature (6)", fmt.Errorf("signature failed on 5")
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return "failed to verify request signature (7)", err
	}

	authToken := ""
	if numKeyParts == 1 {
		account, err := core.ReadAccount(accountID)
		if err != nil {
			return "failed to verify request signature (8)", err
		}
		authToken = account.AuthToken
	} else {
		applicationID, err := strconv.ParseInt(ids[1], 10, 64)
		if err != nil {
			return "failed to verify request signature (9)", err
		}

		application, err := core.ReadApplication(accountID, applicationID)
		if err != nil {
			return "failed to verify request signature (10)", err
		}

		authToken = application.AuthToken
	}

	signString := generateSigningString(ctx.Scope, ctx.Version, ctx.R)

	signingKey := generateSigningKey(authToken, ctx.Scope, ctx.Version, ctx.R)

	/* TODO Debug content, don't remove unless you want to redo it later ** /
	fmt.Printf("\nURL %s\n", ctx.R.URL.Path)
	fmt.Printf("\nPayload %s - %s \n", ctx.R.Header.Get("x-tapglue-payload-hash"), Base64Encode(Sha256String(ctx.BodyString)))
	fmt.Printf("\nSession %s\n", ctx.R.Header.Get("x-tapglue-session"))
	fmt.Printf("\nSignature parts %s - %s \n", signingKey, signString)
	fmt.Printf("\nSignature %s - %s \n\n", ctx.R.Header.Get("x-tapglue-signature"), Base64Encode(Sha256String(signingKey+signString)))
	/**/

	if ctx.R.Header.Get("x-tapglue-signature") == Base64Encode(Sha256String(signingKey+signString)) {
		return "", nil
	}

	return "failed to verify request signature (11)",
		fmt.Errorf("expected %s got %s", ctx.R.Header.Get("x-tapglue-signature"), Base64Encode(Sha256String(signingKey+signString)))
}
