/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import . "gopkg.in/check.v1"

// AddUserSession test to write empty entity
func (dbs *DatabaseSuite) TestAddUserSession_Empty(c *C) {
	// Write session
	savedSession, err := AddUserSession(emptySession)

	// Perform tests
	c.Assert(savedSession, IsNil)
	c.Assert(err, NotNil)
}

// AddUserSession test to write correct session
func (dbs *DatabaseSuite) TestAddUserSession_Correct(c *C) {
	// Prepare data
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

// GetSessionByID test to get an session by its id
func (dbs *DatabaseSuite) TestGetSessionByID_Correct(c *C) {
	// Write correct session
	savedSession := AddCorrectUserSession()

	// Get session by id
	getSession, err := GetSessionByID(savedSession.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getSession, DeepEquals, savedSession)
}

// GetAllUserSessions test to get all users sessions
func (dbs *DatabaseSuite) TestGetAllUserSessions_Correct(c *C) {
	// Write correct sessions
	savedSession1, savedSession2 := AddCorrectUserSessions()

	// Perform tests
	c.Assert(savedSession1, NotNil)
	c.Assert(savedSession2, NotNil)
	c.Assert(savedSession1.AppID, Equals, savedSession2.AppID)
	c.Assert(savedSession1.UserToken, Equals, savedSession2.UserToken)

	// Get account user by id
	getSessions, err := GetAllUserSessions(savedSession1.AppID, savedSession1.UserToken)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getSessions, HasLen, 2)
}

// BenchmarkAddUserSession executes AddUserSession 1000 times
func (dbs *DatabaseSuite) BenchmarkAddUserSession(c *C) {
	// Loop to create 1000 account users
	for i := 0; i < 1000; i++ {
		_, _ = AddUserSession(correctSession)
	}
}
