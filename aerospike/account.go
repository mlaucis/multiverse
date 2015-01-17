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
	client := Client()

	ops := []*as.Operation{
		as.AddOp(as.NewBin("account_id", 1)),
		as.GetOp(),
	}

	var (
		rec *as.Record
		key *as.Key
	)
	if key, err = as.NewKey("keys", "keys", "account_id"); err != nil {
		return
	}
	rec, err = client.Operate(nil, key, ops...)

	return int64(rec.Bins["account_id"].(int)), err
}

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID int64) (account *entity.Account, err error) {
	account = &entity.Account{}

	client := Client()

	var key *as.Key
	if key, err = as.NewKey("accounts", "acc", accountID); err != nil {
		return
	}

	var rec *as.Record
	policy := as.NewPolicy()
	if rec, err = client.Get(policy, key); err != nil {
		return nil, err
	}

	account.ID = int64(rec.Bins["id"].(int))
	account.Name = rec.Bins["name"].(string)

	if rec.Bins["enabled"].(int) == 1 {
		account.Enabled = true
	} else {
		account.Enabled = false
	}

	return
}

// AddAccount adds a new account to the database and returns the created account or an error
func AddAccount(account *entity.Account, retrieve bool) (*entity.Account, error) {
	// Check if name empty
	if account.Name == "" {
		return nil, fmt.Errorf("account name should not be empty")
	}

	var err error

	if account.ID, err = getNewAccountID(); err != nil {
		return nil, err
	}

	// TODO find a better way to store the buckets / sets (and maybe actually learn what those things are in the first place?)
	var key *as.Key
	key, err = as.NewKey("accounts", "acc", account.ID)

	// define some bins with data
	bins := as.BinMap{
		"id":   account.ID,
		"name": account.Name,
	}
	if account.Enabled {
		bins["enabled"] = 1
	} else {
		bins["enabled"] = 0
	}

	client := Client()
	policy := as.NewWritePolicy(0, 0)
	// write the bins
	if err = client.Put(policy, key, bins); err != nil {
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	// Return account
	return GetAccountByID(account.ID)
}
