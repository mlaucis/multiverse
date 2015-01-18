/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package aerospike

import (
	"fmt"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/tapglue/backend/entity"
)

// getnewAccountID generates a new account ID
func getNewAccountID() (accountID int64, err error) {
	ops := []*as.Operation{
		as.AddOp(as.NewBin("account_id", 1)),
		as.GetOp(),
	}

	var rec *as.Record
	rec, err = OperateString("keys", "keys", "account_id", ops, nil)
	if err != nil {
		return 0, err
	}

	return BinToInt64(rec.Bins["account_id"]), err
}

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID int64) (account *entity.Account, err error) {
	var rec *as.Record

	rec, err = GetByInt64("accounts", "acc", accountID, nil)
	if err != nil {
		return nil, err
	}

	account = &entity.Account{}
	account.ID = accountID
	account.Name = BinToString(rec.Bins["name"])
	account.Enabled = BinToBool(rec.Bins["enabled"])

	return
}

// AddAccount adds a new account to the database and returns the created account or an error
func AddAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	// Check if name empty
	if account.Name == "" {
		return nil, fmt.Errorf("account name should not be empty")
	}

	if account.ID, err = getNewAccountID(); err != nil {
		return nil, err
	}

	// define some bins with data
	bins := []*as.Bin{
		as.NewBin("id", account.ID),
		as.NewBin("name", account.Name),
		as.NewBin("enabled", BoolToBin(account.Enabled)),
	}

	err = PutInt64("accounts", "acc", account.ID, bins, nil)
	if err != nil {
		return
	}

	if !retrieve {
		return account, nil
	}

	// Return account
	return GetAccountByID(account.ID)
}
