/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import . "gopkg.in/check.v1"

// WriteAccountUser test to write correct account user
func (cs *CoreSuite) TestWriteAccountUser_Correct(c *C) {
	// Write correct account
	savedAccount := AddCorrectAccount()

	// Write account user
	savedAccountUser, err := WriteAccountUser(correctAccountUser, false)

	// Perform tests
	c.Assert(savedAccountUser, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccountUser.AccountID, Equals, savedAccount.ID)
	c.Assert(savedAccountUser.Username, Equals, correctAccountUser.Username)
	c.Assert(savedAccountUser.Password, Equals, correctAccountUser.Password)
	c.Assert(savedAccountUser.Email, Equals, correctAccountUser.Email)
	c.Assert(savedAccountUser.Enabled, Equals, true)

	// Test types
	c.Assert(savedAccountUser.AccountID, FitsTypeOf, int64(0))
	c.Assert(savedAccountUser.Username, FitsTypeOf, string(""))
	c.Assert(savedAccountUser.Password, FitsTypeOf, string(""))
	c.Assert(savedAccountUser.Email, FitsTypeOf, string(""))
}

// ReadAccountUser test to get an account by its id
func (cs *CoreSuite) TestReadAccountUser_Correct(c *C) {
	// Write correct account user
	savedAccountUser := AddCorrectAccountUser()

	// Get account user by id
	getAccountUser, err := ReadAccountUser(savedAccountUser.AccountID, savedAccountUser.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUser, DeepEquals, savedAccountUser)
}

// ReadAccountUserList test to get all account users
func (cs *CoreSuite) TestReadAccountUserList_Correct(c *C) {
	// Write correct account users
	savedAccountUser1, savedAccountUser2 := AddCorrectAccountUsers()

	// Perform tests
	c.Assert(savedAccountUser1, NotNil)
	c.Assert(savedAccountUser2, NotNil)
	c.Assert(savedAccountUser1.AccountID, Equals, savedAccountUser2.AccountID)

	// Get account user by id
	getAccountUsers, err := ReadAccountUserList(savedAccountUser1.AccountID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccountUsers, HasLen, 2)
}

// BenchmarkWriteAccountUser executes WriteAccountUser 1000 times
func (cs *CoreSuite) BenchmarkWriteAccountUser(c *C) {
	// Loop to create 1000 account users
	for i := 0; i < 1000; i++ {
		_, _ = WriteAccountUser(correctAccountUser, false)
	}
}
