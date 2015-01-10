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
	// Define data
	var (
		accountID   uint64 = 1
		application        = &entity.Application{}
	)

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
	// Define data
	var (
		application = &entity.Application{
			Name: "Demo App",
			Key:  "imanappkey12345",
		}
	)

	// Write account
	savedAccount := AddCorrectAccount()

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

// GetApplicationByID test to get an application by its id
func (dbs *DatabaseSuite) TestGetApplicationByID_Correct(c *C) {
	// Define data
	var (
		application = &entity.Application{
			Name: "Demo App",
			Key:  "imanappkey12345",
		}
	)

	// Write account
	savedAccount := AddCorrectAccount()

	// Write application
	savedApplication, err := AddAccountApplication(savedAccount.ID, application)

	// Perform tests
	c.Assert(savedApplication, NotNil)
	c.Assert(err, IsNil)

	// Get account user by id
	getAccountApplication, err := GetApplicationByID(savedApplication.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountApplication, DeepEquals, savedApplication)
}

// GetAccountAllApplications test to get all applications
func (dbs *DatabaseSuite) TestGetAccountAllApplications_Correct(c *C) {
	// Define data
	var (
		application = &entity.Application{
			Name: "Demo App",
			Key:  "imanappkey12345",
		}
	)

	// Write account
	savedAccount := AddCorrectAccount()

	// Write application
	savedApplication1, err := AddAccountApplication(savedAccount.ID, application)
	savedApplication2, err := AddAccountApplication(savedAccount.ID, application)

	// Perform tests
	c.Assert(savedApplication1, NotNil)
	c.Assert(savedApplication2, NotNil)
	c.Assert(err, IsNil)

	// Get account user by id
	getApplications, err := GetAccountAllApplications(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getApplications, HasLen, 2)
}
