/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

// Test writeResponse
func (s *ServerSuite) TestWriteResponse(c *C) {
	// Implement test
}

// Test errorHappened
func TesterrorHappened(t *testing.T) {
	// Implement test
}

// Test home
func Testhome(t *testing.T) {
	// Implement test
}

// Test humans
func Testhumans(t *testing.T) {
	// Implement test
}

// Test robots
func Testrobots(t *testing.T) {
	// Implement test
}

// Test GetRouter
func TestGetRouter(t *testing.T) {
	// Implement test
}
