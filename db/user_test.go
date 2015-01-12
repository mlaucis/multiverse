/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import . "gopkg.in/check.v1"

// AddApplicationUser test to write empty entity
func (dbs *DatabaseSuite) TestAddApplicationUser_Empty(c *C) {
	// Write application
	savedApplication := AddCorrectAccountApplication()

	// Write user
	savedUser, err := AddApplicationUser(savedApplication.ID, emptyUser)

	// Perform tests
	c.Assert(savedUser, IsNil)
	c.Assert(err, NotNil)
}

// TestAddApplicationUser test to write entity without app
func (dbs *DatabaseSuite) TestAddApplicationUser_NoApp(c *C) {
	// Write user
	savedUser, err := AddApplicationUser(1, correctUser)

	// Perform tests
	c.Assert(savedUser, IsNil)
	c.Assert(err, NotNil)
}

// AddApplicationUser test to write correct user
func (dbs *DatabaseSuite) TestAddApplicationUser_Correct(c *C) {
	// Write application
	savedApplication := AddCorrectAccountApplication()

	// Write user
	savedUser, err := AddApplicationUser(savedApplication.ID, correctUser)

	// Perform tests
	c.Assert(savedUser, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedUser.AppID, Equals, savedApplication.ID)
	c.Assert(savedUser.Token, Equals, correctUser.Token)
	c.Assert(savedUser.Name, Equals, correctUser.Name)
	c.Assert(savedUser.Password, Equals, correctUser.Password)
	c.Assert(savedUser.Email, Equals, correctUser.Email)
	c.Assert(savedUser.ThumbnailURL, Equals, correctUser.ThumbnailURL)
	c.Assert(savedUser.Provider, Equals, correctUser.Provider)
	c.Assert(savedUser.Custom, Equals, correctUser.Custom)
	// Test types
	c.Assert(savedUser.AppID, FitsTypeOf, uint64(0))
	c.Assert(savedUser.Token, FitsTypeOf, string(""))
	c.Assert(savedUser.Name, FitsTypeOf, string(""))
	c.Assert(savedUser.Password, FitsTypeOf, string(""))
	c.Assert(savedUser.Email, FitsTypeOf, string(""))
	c.Assert(savedUser.ThumbnailURL, FitsTypeOf, string(""))
	c.Assert(savedUser.Provider, FitsTypeOf, string(""))
	c.Assert(savedUser.Custom, FitsTypeOf, string(""))
}

// GetApplicationUserByToken test to get user by its token
func (dbs *DatabaseSuite) TestGetApplicationUserByToken_Correct(c *C) {
	// Write correct user
	savedUser := AddCorrectApplicationUser()

	// Get account user by id
	getApplicationUser, err := GetApplicationUserByToken(savedUser.AppID, savedUser.Token)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getApplicationUser, DeepEquals, savedUser)
}

// GetApplicationUsers test to get all users
func (dbs *DatabaseSuite) TestGetApplicationUsers_Correct(c *C) {
	// Write correct users
	savedUser1, savedUser2 := AddCorrectApplicationUsers()

	// Perform tests
	c.Assert(savedUser1, NotNil)
	c.Assert(savedUser2, NotNil)
	c.Assert(savedUser1.AppID, Equals, savedUser2.AppID)

	// Get account user by id
	getUsers, err := GetApplicationUsers(savedUser1.AppID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getUsers, HasLen, 2)
}

// BenchmarkAddApplicationUser executes AddApplicationUser 1000 times
func (dbs *DatabaseSuite) BenchmarkAddApplicationUser(c *C) {
	// Loop to create 1000 account users
	for i := 0; i < 1000; i++ {
		_, _ = AddApplicationUser(1, correctUser)
	}
}
