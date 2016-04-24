package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// UserController bundles the business constraints of Users.
type UserController struct {
	connections connection.Service
	users       user.StrangleService
}

// NewUserController returns a controller instance.
func NewUserController(
	connections connection.Service,
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
) (user.StrangleList, error) {
	us, errs := c.users.FilterByEmail(app.OrgID, app.ID, emails)
	if errs != nil {
		return nil, errs[0]
	}

	for _, u := range us {
		r, err := queryRelation(c.connections, app, originID, u.ID)
		if err != nil {
			return nil, err
		}

		u.Relation = v04_entity.Relation{
			IsFriend:   &r.isFriend,
			IsFollower: &r.isFollower,
			IsFollowed: &r.isFollowing,
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
) (user.StrangleList, error) {
	us, errs := c.users.FilterBySocialIDs(app.OrgID, app.ID, platform, ids)
	if errs != nil {
		return nil, errs[0]
	}

	for _, user := range us {
		r, err := queryRelation(c.connections, app, originID, user.ID)
		if err != nil {
			return nil, err
		}

		user.Relation = v04_entity.Relation{
			IsFriend:   &r.isFriend,
			IsFollower: &r.isFollower,
			IsFollowed: &r.isFollowing,
		}
	}

	return us, nil
}
