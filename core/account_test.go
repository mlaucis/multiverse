/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"
	// FIX import cycle not allowed in test

	. "gopkg.in/check.v1"
)

// WriteAccount test to write account entity with just a name
func (cs *CoreSuite) TestWriteAccount_Correct(c *C) {
	// Write account
	savedAccount, err := utils.AddCorrectAccount()

	// Perform tests
	c.Assert(savedAccount, NotNil)
	c.Assert(err, IsNil)
	c.Assert(savedAccount.Enabled, Equals, true)
}

// ReadAccount test to get an account by its id
func (cs *CoreSuite) TestReadAccount(c *C) {
	// Write correct account
	savedAccount, err := utils.AddCorrectAccount()

	// Get account by id
	getAccount, err := ReadAccount(savedAccount.ID)

	// Perform tests
	c.Assert(err, IsNil)
	c.Assert(getAccount, DeepEquals, savedAccount)
}

// BenchmarkAccountStep1Write executes WriteAccount 1000 times
func (cs *CoreSuite) BenchmarkAccountStep1Write(c *C) {
	var i int64

	for i = 1; i <= 1000; i++ {
		correctAccount.ID = i
		_, _ = utils.AddCorrectAccount()
	}
}

// BenchmarkAccountStep2Write executes ReadAccount 1000 times
func (cs *CoreSuite) BenchmarkAccountStep2Read(c *C) {
	var (
		i        int64
		accounts = make(map[int64]*entity.Account)
	)

	for i = 1; i <= 1000; i++ {
		account, _ := ReadAccount(i)
		accounts[i] = account
	}

	accounts = nil
}
