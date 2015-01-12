/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import . "gopkg.in/check.v1"

// AddUserSession test to write empty entity
func (dbs *DatabaseSuite) TestAddUserSession_Empty(c *C) {
	// Write session
	savedSession, err := AddUserSession(emptySessions)

	// Perform tests
	c.Assert(savedSession, IsNil)
	c.Assert(err, NotNil)
}

// AddUserSession test to write correct session
func (dbs *DatabaseSuite) TestAddUserSession_Correct(c *C) {
	// Prepare data
	correctUser.Token = RandomToken()
	savedUser := AddCorrectApplicationUser()
	correctSession.UserToken = savedUser.Token
	correctSession.AppID = savedUser.AppID

	// Write session
	savedSession, err := AddUserSession(correctSession)

	// Perform tests
	c.Assert(savedSession, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedSession.AppID, Equals, savedUser.AppID)
	c.Assert(savedSession.UserToken, Equals, savedUser.Token)
	// Test types
	c.Assert(savedSession.AppID, FitsTypeOf, uint64(0))
	c.Assert(savedSession.UserToken, FitsTypeOf, string(""))
	c.Assert(savedSession.Custom, FitsTypeOf, string(""))
	c.Assert(savedUser.Password, FitsTypeOf, string(""))
	c.Assert(savedUser.Email, FitsTypeOf, string(""))
	c.Assert(savedUser.ThumbnailURL, FitsTypeOf, string(""))
	c.Assert(savedUser.Provider, FitsTypeOf, string(""))
	c.Assert(savedUser.Custom, FitsTypeOf, string(""))
}

// BenchmarkAddUserSession executes AddUserSession 1000 times
func (dbs *DatabaseSuite) BenchmarkAddUserSession(c *C) {
	// Loop to create 1000 account users
	for i := 0; i < 1000; i++ {
		_, _ = AddUserSession(correctSession)
	}
}
