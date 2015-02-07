/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package mysql

import . "gopkg.in/check.v1"

// AddSessionEvent test to write empty entity
func (dbs *DatabaseSuite) TestAddSessionEvent_Empty(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write session
		savedSession := AddCorrectUserSession()
		UpdateEvent(savedSession.applicationId, savedSession.ID, savedSession.UserToken)

		// Write event
		savedEvent, err := AddSessionEvent(emptyEvent)

		// Perform tests
		c.Assert(savedEvent, IsNil)
		c.Assert(err, NotNil)
	*/
}

// AddSessionEvent test to write correct event
func (dbs *DatabaseSuite) TestAddSessionEvent_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Prepare data
		savedSession := AddCorrectUserSession()
		UpdateEvent(savedSession.applicationId, savedSession.ID, savedSession.UserToken)

		// Write event
		savedEvent, err := AddSessionEvent(correctEvent)

		// Perform tests
		c.Assert(err, IsNil)
		c.Assert(savedEvent, NotNil)
		c.Assert(savedEvent.ApplicationID, Equals, savedSession.ApplicationID)
		c.Assert(savedEvent.UserID, Equals, savedSession.UserID)
		// Test types
		c.Assert(savedEvent.ApplicationID, FitsTypeOf, uint64(0))
		c.Assert(savedEvent.UserID, FitsTypeOf, uint64(0))
		c.Assert(savedEvent.Metadata, FitsTypeOf, string(""))
		c.Assert(savedEvent.Object.DisplayName["en"], FitsTypeOf, string(""))
		c.Assert(savedEvent.Metadata, FitsTypeOf, string(""))
	*/
}

// GetEventByID test to get event by its id
func (dbs *DatabaseSuite) TestGetEventByID_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct event
		savedEvent := AddCorrectEvent()

		// Get event user by id
		getEvent, err := GetEventByID(savedEvent.ID)

		// Perform tests
		c.Assert(err, IsNil)
		c.Assert(getEvent, DeepEquals, savedEvent)
	*/
}

// GetAllUserAppEvents test to get all users events
func (dbs *DatabaseSuite) TestGetAllUserAppEvents_Correct(c *C) {
	c.Skip("not refactored yet")
	/*
		// Write correct sessions
		savedEvent1, savedEvent2 := AddCorrectEvents()

		// Perform tests
		c.Assert(savedEvent1, NotNil)
		c.Assert(savedEvent2, NotNil)
		c.Assert(savedEvent1.ApplicationID, Equals, savedEvent2.ApplicationID)
		c.Assert(savedEvent1.UserID, Equals, savedEvent2.UserID)

		// Get account user by id
		getEvents, err := GetAllUserAppEvents(savedEvent1.ApplicationID, savedEvent1.UserID)

		// Perform tests
		c.Assert(err, IsNil)
		c.Assert(getEvents, NotNil)
	*/
}

// BenchmarkAddSessionEvent executes AddSessionEvent 1000 times
func (dbs *DatabaseSuite) BenchmarkAddSessionEvent(c *C) {
	c.Skip("not refactored yet")
	/*
		// Loop to create 1000 events
		for i := 0; i < 1000; i++ {
			_, _ = AddSessionEvent(correctEvent)
		}
	*/
}
