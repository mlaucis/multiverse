/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import . "gopkg.in/check.v1"

// AddSessionEvent test to write empty entity
func (dbs *DatabaseSuite) TestAddSessionEvent_Empty(c *C) {
	// Write session
	savedSession := AddCorrectUserSession()
	UpdateEvent(savedSession.AppID, savedSession.ID, savedSession.UserToken)

	// Write event
	savedEvent, err := AddSessionEvent(emptyEvent)

	// Perform tests
	c.Assert(savedEvent, IsNil)
	c.Assert(err, NotNil)
}

// AddSessionEvent test to write correct event
func (dbs *DatabaseSuite) TestAddSessionEvent_Correct(c *C) {
	// Prepare data
	savedSession := AddCorrectUserSession()
	UpdateEvent(savedSession.AppID, savedSession.ID, savedSession.UserToken)

	// Write event
	savedEvent, err := AddSessionEvent(correctEvent)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(savedEvent, NotNil)
	c.Assert(savedEvent.AppID, Equals, savedSession.AppID)
	c.Assert(savedEvent.UserToken, Equals, savedSession.UserToken)
	c.Assert(savedEvent.SessionID, Equals, savedSession.ID)
	// Test types
	c.Assert(savedEvent.AppID, FitsTypeOf, uint64(0))
	c.Assert(savedEvent.SessionID, FitsTypeOf, uint64(0))
	c.Assert(savedEvent.UserToken, FitsTypeOf, string(""))
	c.Assert(savedEvent.Custom, FitsTypeOf, string(""))
	c.Assert(savedEvent.Title, FitsTypeOf, string(""))
	c.Assert(savedEvent.Custom, FitsTypeOf, string(""))
}

// GetEventByID test to get event by its id
func (dbs *DatabaseSuite) TestGetEventByID_Correct(c *C) {
	// Write correct event
	savedEvent := AddCorrectEvent()

	// Get event user by id
	getEvent, err := GetEventByID(savedEvent.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getEvent, DeepEquals, savedEvent)
}

// GetAllUserAppEvents test to get all users events
func (dbs *DatabaseSuite) TestGetAllUserAppEvents_Correct(c *C) {
	// Write correct sessions
	savedEvent1, savedEvent2 := AddCorrectEvents()

	// Perform tests
	c.Assert(savedEvent1, NotNil)
	c.Assert(savedEvent2, NotNil)
	c.Assert(savedEvent1.AppID, Equals, savedEvent2.AppID)
	c.Assert(savedEvent1.UserToken, Equals, savedEvent2.UserToken)
	c.Assert(savedEvent1.SessionID, Equals, savedEvent2.SessionID)

	// Get account user by id
	getEvents, err := GetAllUserAppEvents(savedEvent1.AppID, savedEvent1.UserToken)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getEvents, NotNil)
}

// GetSessionEvents test to get all session events
func (dbs *DatabaseSuite) TestGetGetSessionEvents_Correct(c *C) {
	// Write correct sessions
	savedEvent1, savedEvent2 := AddCorrectEvents()

	// Perform tests
	c.Assert(savedEvent1, NotNil)
	c.Assert(savedEvent2, NotNil)
	c.Assert(savedEvent1.AppID, Equals, savedEvent2.AppID)
	c.Assert(savedEvent1.UserToken, Equals, savedEvent2.UserToken)
	c.Assert(savedEvent1.SessionID, Equals, savedEvent2.SessionID)

	// Get account user by id
	getEvents, err := GetSessionEvents(savedEvent1.AppID, savedEvent1.SessionID, savedEvent1.UserToken)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getEvents, NotNil)
}

// BenchmarkAddSessionEvent executes AddSessionEvent 1000 times
func (dbs *DatabaseSuite) BenchmarkAddSessionEvent(c *C) {
	// Loop to create 1000 events
	for i := 0; i < 1000; i++ {
		_, _ = AddSessionEvent(correctEvent)
	}
}
