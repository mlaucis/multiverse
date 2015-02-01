/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package mysql

import . "gopkg.in/check.v1"

// AddAccountApplication test to write empty entity
func (dbs *DatabaseSuite) TestAddAccountApplication_Empty(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct account
		savedAccount := AddCorrectAccount()

		// Write application
		savedApplication, err := AddAccountApplication(savedAccount.ID, emptyApplication)

		// Perform tests
		c.Assert(savedApplication, IsNil)
		c.Assert(err, NotNil)
	*/
}

// AddAccountApplication test to write an application without an account
func (dbs *DatabaseSuite) TestAddAccountApplication_NoAccount(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write application
		savedApplication, err := AddAccountApplication(0, correctApplication)

		// Perform tests
		c.Assert(savedApplication, IsNil)
		c.Assert(err, NotNil)
	*/
}

// AddAccountApplication test to write application entity with name and key
func (dbs *DatabaseSuite) TestAddAccountApplication_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct account
		savedAccount := AddCorrectAccount()

		// Write application
		savedApplication, err := AddAccountApplication(savedAccount.ID, correctApplication)

		// Perform tests
		c.Assert(savedApplication, NotNil)
		c.Assert(err, IsNil)
		c.Assert(savedApplication.AccountID, Equals, savedAccount.ID)
		c.Assert(savedApplication.Name, Equals, correctApplication.Name)
		c.Assert(savedApplication.AuthToken, Equals, correctApplication.AuthToken)
		c.Assert(savedApplication.Enabled, Equals, true)
		// Test types
		c.Assert(savedApplication.AuthToken, FitsTypeOf, string(""))
		c.Assert(savedApplication.Name, FitsTypeOf, string(""))
		c.Assert(savedApplication.AccountID, FitsTypeOf, uint64(0))
		c.Assert(savedApplication.Enabled, FitsTypeOf, bool(true))
	*/
}

// GetApplicationByID test to get an application by its id
func (dbs *DatabaseSuite) TestGetApplicationByID_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct application
		savedApplication := AddCorrectAccountApplication()

		// Get account user by id
		getAccountApplication, err := GetApplicationByID(savedApplication.ID)

		// Perform tests
		c.Assert(err, IsNil)
		c.Assert(getAccountApplication, DeepEquals, savedApplication)
	*/
}

// GetAccountAllApplications test to get all applications
func (dbs *DatabaseSuite) TestGetAccountAllApplications_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct applications
		savedApplication1, savedApplication2 := AddCorrectAccountApplications()

		// Perform tests
		c.Assert(savedApplication1, NotNil)
		c.Assert(savedApplication2, NotNil)
		c.Assert(savedApplication1.AccountID, Equals, savedApplication2.AccountID)

		// Get account user by id
		getApplications, err := GetAccountAllApplications(savedApplication1.AccountID)

		// Perform tests
		c.Assert(err, IsNil)
		c.Assert(getApplications, HasLen, 2)
	*/
}

// BenchmarkAddAccountApplication executes ddAccountApplication 1000 times
func (dbs *DatabaseSuite) BenchmarkAddAccountApplication(c *C) {
	c.Skip("not refactored yet")
	/*
		// Loop to create 1000 applications
		for i := 0; i < 1000; i++ {
			_, _ = AddAccountApplication(1, correctApplication)
		}
	*/
}
