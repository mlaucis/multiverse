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
	InitDatabases(cfg.DB())

	var account = &entity.Account{}

	savedAccount, err := AddAccount(account)

	c.Assert(savedAccount, IsNil)
	c.Assert(err, NotNil)
}

// AddAccount test to write account entity with just a name
func (dbs *DatabaseSuite) TestAddAccount_Correct(c *C) {
	InitDatabases(cfg.DB())

	var account = &entity.Account{
		Name: "Demo",
	}

	savedAccount, err := AddAccount(account)

	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Name, Equals, account.Name)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// Test GetAccountByID
func (dbs *DatabaseSuite) TestGetAccountByID(c *C) {
	InitDatabases(cfg.DB())

	var account = &entity.Account{
		Name: "Demo",
	}

	// Write account first
	savedAccount, err := AddAccount(account)

	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)

	// Get account by id
	getAccount, err := GetAccountByID(savedAccount.ID)

	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}
