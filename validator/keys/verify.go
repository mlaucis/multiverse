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

// VerifyRequest verifies if a request is properly signed or not
func VerifyRequest(ctx *context.Context, numKeyParts int) error {
	signature := ctx.R.Header.Get("x-tapglue-signature")
	if signature == "" {
		return fmt.Errorf("signature failed on 1")
	}

	if _, err := Base64Decode(signature); err != nil {
		return fmt.Errorf("signature failed on 2")
	}

	if Base64Encode(Sha256String(ctx.BodyString)) != ctx.R.Header.Get("x-tapglue-payload-hash") {
		return fmt.Errorf("signature failed on 3")
	}

	encodedIds := ctx.R.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return fmt.Errorf("signature failed on 4")
	}

	ids := strings.SplitN(string(decodedIds), ":", numKeyParts)
	if len(ids) != numKeyParts {
		return fmt.Errorf("signature failed on 5")
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return fmt.Errorf("signature failed on 6")
	}

	authToken := ""
	if numKeyParts == 1 {
		account, err := core.ReadAccount(accountID)
		if err != nil {
			return fmt.Errorf("signature failed on 7")
		}
		authToken = account.AuthToken
	} else {
		applicationID, err := strconv.ParseInt(ids[1], 10, 64)
		if err != nil {
			return fmt.Errorf("signature failed on 8")
		}

		application, err := core.ReadApplication(accountID, applicationID)
		if err != nil {
			return fmt.Errorf("signature failed on 9")
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

	if ctx.R.Header.Get("x-tapglue-signature") != Base64Encode(Sha256String(signingKey+signString)) {
		return fmt.Errorf("signature failed on 10")
	}
	return nil
}
