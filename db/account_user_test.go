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
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
	var (
		accountID   uint64
		accountUser = &entity.AccountUser{}
	)

	accountID = 1

	// Write account user
	savedAccountUser, err := AddAccountUser(accountID, accountUser)

	// Perform tests
	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}

// AddAccountUser test to write account user entity without an account
func (dbs *DatabaseSuite) TestAddAccountUser_NoAccount(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
	var (
		accountID   uint64
		accountUser = &entity.AccountUser{
			Name:     "Demo User",
			Password: "iamsecure..not",
			Email:    "d@m.o",
		}
	)

	accountID = 0

	// Write account user
	savedAccountUser, err := AddAccountUser(accountID, accountUser)

	// Perform tests
	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}

// AddAccountUser test to write account user entity with name, pw, email
func (dbs *DatabaseSuite) TestAddAccountUser_Correct(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
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

	// Write account
	savedAccount, err := AddAccount(account)

	// Write account user
	savedAccountUser, err := AddAccountUser(savedAccount.ID, accountUser)

	// Perform tests
	c.Assert(savedAccountUser, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccountUser.AccountID, Equals, savedAccount.ID)
	c.Assert(savedAccountUser.Name, Equals, accountUser.Name)
	c.Assert(savedAccountUser.Password, Equals, accountUser.Password)
	c.Assert(savedAccountUser.Email, Equals, accountUser.Email)
	c.Assert(savedAccountUser.Enabled, Equals, true)
}

// GetAccountUserByID test to get an account by its id
func (dbs *DatabaseSuite) TestGetAccountUserByID_Correct(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
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

	// Write account
	savedAccount, err := AddAccount(account)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)

	// Write account user
	savedAccountUser, err := AddAccountUser(savedAccount.ID, accountUser)

	// Perform tests
	c.Assert(savedAccountUser, NotNil)
	c.Assert(err, IsNil)

	// Get account user by id
	getAccountUser, err := GetAccountUserByID(savedAccount.ID, savedAccountUser.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUser, DeepEquals, savedAccountUser)
}

// GetAccountAllUsers test to get all account users
func (dbs *DatabaseSuite) TestGetAccountAllUsers_Correct(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
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

	// Write account
	savedAccount, err := AddAccount(account)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)

	// Write account user
	savedAccountUser1, err := AddAccountUser(savedAccount.ID, accountUser)

	// Write another account user
	savedAccountUser2, err := AddAccountUser(savedAccount.ID, accountUser)

	// Perform tests
	c.Assert(savedAccountUser1, NotNil)
	c.Assert(savedAccountUser2, NotNil)
	c.Assert(err, IsNil)

	// Get account user by id
	getAccountUsers, err := GetAccountAllUsers(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUsers, HasLen, 2)
}
