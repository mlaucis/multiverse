/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import . "gopkg.in/check.v1"

// AddAccountUser test to write empty entity
func (dbs *DatabaseSuite) TestAddAccountUser_Empty(c *C) {
	// Write correct account
	savedAccount := AddCorrectAccount()

	// Write account user
	savedAccountUser, err := AddAccountUser(savedAccount.ID, emtpyAccountUser)

	// Perform tests
	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}

// AddAccountUser test to write entity without account
func (dbs *DatabaseSuite) TestAddAccountUser_NoAccount(c *C) {
	// Write account user
	savedAccountUser, err := AddAccountUser(0, correctAccountUser)

	// Perform tests
	c.Assert(savedAccountUser, IsNil)
	c.Assert(err, NotNil)
}

// AddAccountUser test to write correct account user
func (dbs *DatabaseSuite) TestAddAccountUser_Correct(c *C) {
	// Write correct account
	savedAccount := AddCorrectAccount()

	// Write account user
	savedAccountUser, err := AddAccountUser(savedAccount.ID, correctAccountUser)

	// Perform tests
	c.Assert(savedAccountUser, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccountUser.AccountID, Equals, savedAccount.ID)
	c.Assert(savedAccountUser.Name, Equals, correctAccountUser.Name)
	c.Assert(savedAccountUser.Password, Equals, correctAccountUser.Password)
	c.Assert(savedAccountUser.Email, Equals, correctAccountUser.Email)
	c.Assert(savedAccountUser.Enabled, Equals, true)
	// Test types
	c.Assert(savedAccountUser.AccountID, FitsTypeOf, uint64(0))
	c.Assert(savedAccountUser.Name, FitsTypeOf, string(""))
	c.Assert(savedAccountUser.Name, FitsTypeOf, string(""))
	c.Assert(savedAccountUser.Password, FitsTypeOf, string(""))
	c.Assert(savedAccountUser.Email, FitsTypeOf, string(""))
}

// GetAccountUserByID test to get an account by its id
func (dbs *DatabaseSuite) TestGetAccountUserByID_Correct(c *C) {
	// Write correct account user
	savedAccountUser := AddCorrectAccountUser()

	// Get account user by id
	getAccountUser, err := GetAccountUserByID(savedAccountUser.AccountID, savedAccountUser.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUser, DeepEquals, savedAccountUser)
}

// GetAccountAllUsers test to get all account users
func (dbs *DatabaseSuite) TestGetAccountAllUsers_Correct(c *C) {
	// Write correct account users
	savedAccountUser1, savedAccountUser2 := AddCorrectAccountUsers()

	// Perform tests
	c.Assert(savedAccountUser1, NotNil)
	c.Assert(savedAccountUser2, NotNil)
	c.Assert(savedAccountUser1.AccountID, Equals, savedAccountUser2.AccountID)

	// Get account user by id
	getAccountUsers, err := GetAccountAllUsers(savedAccountUser1.AccountID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUsers, HasLen, 2)
}

// BenchmarkAddAccountUser executes AddAccountUser 1000 times
func (dbs *DatabaseSuite) BenchmarkAddAccountUser(c *C) {
	// Loop to create 1000 account users
	for i := 0; i < 1000; i++ {
		_, _ = AddAccountUser(1, correctAccountUser)
	}
}
