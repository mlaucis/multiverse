/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import (
	"github.com/tapglue/backend/entity"

	. "gopkg.in/check.v1"
)

// AddAccountUser test to write empty entity
func (dbs *DatabaseSuite) TestAddAccountUser_Empty(c *C) {
	InitDatabases(cfg.DB())

	var (
		accountID   uint64
		accountUser = &entity.AccountUser{}
	)

	accountID = 1

	savedAccountUser, err := AddAccountUser(accountID, accountUser)

	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}

// AddAccountUser test to write account user entity with name, pw, email
func (dbs *DatabaseSuite) TestAddAccountUser_Correct(c *C) {
	InitDatabases(cfg.DB())

	var (
		account = &entity.Account{
			Name: "Demo",
		}
		accountUser = &entity.AccountUser{
			Name:     "Demo User",
			Password: "iamsecure..not",
			Email:    "d@m.o",
		}
	)

	// Write account first
	savedAccount, err := AddAccount(account)

	savedAccountUser, err := AddAccountUser(savedAccount.ID, accountUser)

	c.Assert(savedAccountUser, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccountUser.AccountID, Equals, savedAccount.ID)
	c.Assert(savedAccountUser.Name, Equals, accountUser.Name)
	c.Assert(savedAccountUser.Password, Equals, accountUser.Password)
	c.Assert(savedAccountUser.Email, Equals, accountUser.Email)
	c.Assert(savedAccountUser.Enabled, Equals, true)
}

// AddAccountUser test to write account user entity without an account
func (dbs *DatabaseSuite) TestAddAccountUser_NoAccount(c *C) {
	InitDatabases(cfg.DB())

	var (
		accountID   uint64
		accountUser = &entity.AccountUser{
			Name:     "Demo User",
			Password: "iamsecure..not",
			Email:    "d@m.o",
		}
	)

	accountID = 2

	savedAccountUser, err := AddAccountUser(accountID, accountUser)

	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}
