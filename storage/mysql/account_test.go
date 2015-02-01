/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package mysql

import . "gopkg.in/check.v1"

// AddAccount test to write empty entity
func (dbs *DatabaseSuite) TestAddAccount_Empty(c *C) {
	// Write account
	savedAccount, err := AddAccount(emtpyAccount)

	// Perform tests
	c.Assert(savedAccount, IsNil)
	c.Assert(err, NotNil)
}

// AddAccount test to write account entity with just a name
func (dbs *DatabaseSuite) TestAddAccount_Correct(c *C) {
	c.Skip("not refactored yet")
	// Write account
	savedAccount, err := AddAccount(correctAccount)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Name, Equals, correctAccount.Name)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// GetAccountByID test to get an account by its id
func (dbs *DatabaseSuite) TestGetAccountByID(c *C) {
	// Write correct account
	savedAccount := AddCorrectAccount()

	// Get account by id
	getAccount, err := GetAccountByID(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}

// BenchmarkAddAccount executes AddAccount 1000 times
func (dbs *DatabaseSuite) BenchmarkAddAccount(c *C) {
	// Loop to create 1000 accounts
	for i := 0; i < 1000; i++ {
		_, _ = AddAccount(correctAccount)
	}
}
