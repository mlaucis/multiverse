/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import (
	"github.com/tapglue/backend/entity"

	. "gopkg.in/check.v1"
)

// AddAccountApplication test to write empty entity
func (dbs *DatabaseSuite) TestAddAccountApplication_Empty(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
	var (
		accountID   uint64
		application = &entity.Application{}
	)

	accountID = 1

	// Write application
	savedApplication, err := AddAccountApplication(accountID, application)

	// Perform tests
	c.Assert(savedApplication, IsNil)
	c.Assert(err, NotNil)

	// TBD define expected error
	// c.Assert(err, ErrorMatches, "expected.*error")
}

// AddAccountApplication test to write an application without an account
func (dbs *DatabaseSuite) TestAddAccountApplication_NoAccount(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
	var (
		accountID   uint64
		application = &entity.Application{
			Name: "Demo App",
			Key:  "imanappkey12345",
		}
	)

	accountID = 0

	// Write application
	savedApplication, err := AddAccountApplication(accountID, application)

	// Perform tests
	c.Assert(savedApplication, IsNil)
	c.Assert(err, NotNil)

	// TBD define expected error
	// c.Assert(err, ErrorMatches, "expected.*error")
}

// AddAccountApplication test to write application entity with name and key
func (dbs *DatabaseSuite) TestAddAccountApplication_Correct(c *C) {
	// Initialize database
	InitDatabases(cfg.DB())

	// Define data
	var (
		account = &entity.Account{
			Name: "Demo",
		}
		application = &entity.Application{
			Name: "Demo App",
			Key:  "imanappkey12345",
		}
	)

	// Write account first
	savedAccount, err := AddAccount(account)

	// Write application
	savedApplication, err := AddAccountApplication(savedAccount.ID, application)

	// Perform tests
	c.Assert(savedApplication, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedApplication.AccountID, Equals, savedAccount.ID)
	c.Assert(savedApplication.AccountID, FitsTypeOf, uint64(0))
	c.Assert(savedApplication.Name, Equals, application.Name)
	c.Assert(savedApplication.Name, FitsTypeOf, string("go"))
	c.Assert(savedApplication.Key, Equals, application.Key)
	c.Assert(savedApplication.Key, FitsTypeOf, string("go"))
	c.Assert(savedApplication.Enabled, Equals, true)
	c.Assert(savedApplication.Enabled, FitsTypeOf, bool(true))
}
