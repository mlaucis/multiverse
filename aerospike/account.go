/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package aerospike

import (
	"fmt"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/tapglue/backend/entity"
)

var ass *as.Client

func aerospiked() (client *as.Client) {
	if ass != nil {
		return ass
	}

	var err error
	if client, err = as.NewClient("127.0.0.1", 3000); err != nil {
		panic(err)
	}

	ass = client

	return
}

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID int64) (account *entity.Account, err error) {
	account = &entity.Account{}

	client := aerospiked()

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

	// TODO find a better way to store the buckets / sets (and maybe actually learn what those things are in the first place?)
	key, err := as.NewKey("accounts", "acc", account.ID)

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

	client := aerospiked()
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
