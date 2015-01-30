/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import . "gopkg.in/check.v1"

// WriteAccount test to write empty entity
func (cs *CoreSuite) TestWriteAccount_Empty(c *C) {
	// Write account
	savedAccount, err := WriteAccount(emtpyAccount, true)

	// Perform tests
	c.Assert(savedAccount, IsNil)
	c.Assert(err, NotNil)
}

// WriteAccount test to write account entity with just a name
func (cs *CoreSuite) TestWriteAccount_Correct(c *C) {
	// Write account
	savedAccount, err := WriteAccount(correctAccount, true)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Name, Equals, correctAccount.Name)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// ReadAccount test to get an account by its id
func (cs *CoreSuite) TestReadAccount(c *C) {
	// Write correct account
	savedAccount := AddCorrectAccount()

	// Get account by id
	getAccount, err := ReadAccount(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}

// BenchmarkWriteAccount executes WriteAccount 1000 times
func (cs *CoreSuite) BenchmarkWriteAccount(c *C) {
	var i int64
	// Loop to create 1000 accounts
	for i = 1; i <= 1000; i++ {
		correctAccount.ID = i
		_, _ = WriteAccount(correctAccount, false)
	}
}

// BenchmarkReadAccount executes ReadAccount 1000 times
func (cs *CoreSuite) BenchmarkReadAccount(c *C) {
	var i int64
	// Loop to create 1000 accounts
	for i = 1; i <= 1000; i++ {
		_, _ = ReadAccount(i)
	}
}
