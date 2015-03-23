/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core_test

import (
	"testing"

	. "github.com/tapglue/backend/v01/core"
	. "gopkg.in/check.v1"
)

// WriteAccount test to write account entity with just a name
func (cs *CoreSuite) TestWriteAccount_Correct(c *C) {
	// Write account
	savedAccount, err := AddCorrectAccount(true)

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// ReadAccount test to get an account by its id
func (cs *CoreSuite) TestReadAccount(c *C) {
	// Write correct account
	savedAccount, err := AddCorrectAccount(true)

	// Get account by id
	getAccount, err := ReadAccount(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}

func BenchmarkAccountStep1Write(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		correctAccount.ID = int64(i)
		_, _ = AddCorrectAccount(false)
	}
}

func BenchmarkAccountStep2Read(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		_, _ = ReadAccount(int64(i))
	}
}
