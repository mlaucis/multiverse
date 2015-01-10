/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import (
	"github.com/tapglue/backend/entity"

	. "gopkg.in/check.v1"
)

// AddAccount test to write empty entity
func (dbs *DatabaseSuite) TestAddAccount_Empty(c *C) {
	var account = &entity.Account{}

	// Write account
	savedAccount, err := AddAccount(account)

	// Perform tests
	c.Assert(savedAccount, IsNil)
	c.Assert(err, NotNil)
}

// AddAccount test to write account entity with just a name
func (dbs *DatabaseSuite) TestAddAccount_Correct(c *C) {
	// Define data
	var account = &entity.Account{
		Name: "Demo",
	}

	// Write account
	savedAccount, err := AddAccount(account)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Name, Equals, account.Name)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// GetAccountByID test to get an account by its id
func (dbs *DatabaseSuite) TestGetAccountByID(c *C) {
	// Write account first
	savedAccount := AddCorrectAccount()

	// Get account by id
	getAccount, err := GetAccountByID(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}

func (dbs *DatabaseSuite) BenchmarkAddAccount(c *C) {
	var account = &entity.Account{
		Name: "Demo",
	}

	for i := 0; i < 1000; i++ {
		_, _ = AddAccount(account)
	}
}
