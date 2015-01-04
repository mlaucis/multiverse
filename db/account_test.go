/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package db

import (
	"github.com/tapglue/backend/entity"

	. "gopkg.in/check.v1"
)

// Test GetAccountByID
func (dbs *DatabaseSuite) TestGetAccountByID(c *C) {
	c.Skip("not implemented yet")
}

func (dbs *DatabaseSuite) TestAddAccount_Empty(c *C) {
	InitDatabases(cfg.DB())

	var account = &entity.Account{}

	_, err := AddAccount(account)

	c.Assert(err, Not(IsNil))
}
