package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// UserController bundles the business constraints of Users.
type UserController struct {
	connections connection.StrangleService
	users       user.StrangleService
}

// NewUserController returns a controller instance.
func NewUserController(
	connections connection.StrangleService,
	users user.StrangleService,
) *UserController {
	return &UserController{
		connections: connections,
		users:       users,
	}
}

// ListByEmails returns all users for the given emails.
func (c *UserController) ListByEmails(
	app *v04_entity.Application,
	originID uint64,
	emails ...string,
) (user.List, error) {
	us, errs := c.users.FilterByEmail(app.OrgID, app.ID, emails)
	if errs != nil {
		return nil, errs[0]
	}

	for _, user := range us {
		r, errs := c.connections.Relation(app.OrgID, app.ID, originID, user.ID)
		if errs != nil {
			return nil, errs[0]
		}

		if r != nil {
			user.Relation = *r
		}
	}

	return us, nil
}

// ListByPlatformIDs returns all users for the given ids for the social platform.
func (c *UserController) ListByPlatformIDs(
	app *v04_entity.Application,
	originID uint64,
	platform string,
	ids ...string,
) (user.List, error) {
	us, errs := c.users.FilterBySocialIDs(app.OrgID, app.ID, platform, ids)
	if errs != nil {
		return nil, errs[0]
	}

	for _, user := range us {
		r, errs := c.connections.Relation(app.OrgID, app.ID, originID, user.ID)
		if errs != nil {
			return nil, errs[0]
		}

		if r != nil {
			user.Relation = *r
		}
	}

	return us, nil
}
