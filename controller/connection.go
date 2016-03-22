package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	errmsg "github.com/tapglue/multiverse/v04/errmsg"
)

// ConnectionFeed is the composite to transport information relevant for
// connections.
type ConnectionFeed struct {
	Connections connection.List
	UserMap     user.Map
}

// ConnectionController bundles the business constraints of Connections.
type ConnectionController struct {
	connections connection.Service
	users       user.StrangleService
}

// NewConnectionController returns a controller instance.
func NewConnectionController(
	connections connection.Service,
	users user.StrangleService,
) *ConnectionController {
	return &ConnectionController{
		connections: connections,
		users:       users,
	}
}

// ByState returns all connections for the given origin and state.
func (c *ConnectionController) ByState(
	app *v04_entity.Application,
	originID uint64,
	state connection.State,
) (*ConnectionFeed, error) {
	switch state {
	case connection.StatePending, connection.StateConfirmed, connection.StateRejected:
		// valid
	default:
		return nil, wrapError(ErrInvalidEntity, "unsupported state %s", string(state))
	}

	ics, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{originID},
		States:  []connection.State{state},
	})
	if err != nil {
		return nil, err
	}

	ocs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		States:  []connection.State{state},
		ToIDs:   []uint64{originID},
	})
	if err != nil {
		return nil, err
	}

	um, err := user.MapFromIDs(c.users, app, append(ics.ToIDs(), ocs.FromIDs()...)...)
	if err != nil {
		return nil, err
	}

	return &ConnectionFeed{
		Connections: append(ics, ocs...),
		UserMap:     um,
	}, nil
}

// CreateSocial connects the origin with the users matching the platform ids.
func (c *ConnectionController) CreateSocial(
	app *v04_entity.Application,
	originID uint64,
	connectionType connection.Type,
	connectionState connection.State,
	platform string,
	connectionIDs ...string,
) (user.List, error) {
	us, errs := c.users.FilterBySocialIDs(app.OrgID, app.ID, platform, connectionIDs)
	if errs != nil {
		return nil, errs[0]
	}

	for _, u := range us {
		_, err := c.connections.Put(app.Namespace(), &connection.Connection{
			Enabled: true,
			FromID:  originID,
			ToID:    u.ID,
			State:   connectionState,
			Type:    connectionType,
		})
		if err != nil {
			return nil, err
		}

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

// Delete disables the given connection.
func (c *ConnectionController) Delete(
	app *v04_entity.Application,
	con *connection.Connection,
) error {
	var (
		fromIDs = []uint64{con.FromID}
		toIDs   = []uint64{con.ToID}
	)

	if con.Type == connection.TypeFriend {
		fromIDs = []uint64{con.FromID, con.ToID}
		toIDs = []uint64{con.FromID, con.ToID}
	}

	cs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: fromIDs,
		ToIDs:   toIDs,
		Types:   []connection.Type{con.Type},
	})
	if err != nil {
		return err
	}

	if len(cs) == 0 {
		return nil
	}

	con = cs[0]

	con.Enabled = false

	_, err = c.connections.Put(app.Namespace(), con)

	return err
}

// Followers returns the list of users who follow the origin.
func (c *ConnectionController) Followers(
	app *v04_entity.Application,
	originID uint64,
) (user.List, error) {
	cs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		ToIDs:   []uint64{originID},
		States:  []connection.State{connection.StateConfirmed},
		Types:   []connection.Type{connection.TypeFollow},
	})
	if err != nil {
		return nil, err
	}

	return user.ListFromIDs(c.users, app, cs.FromIDs()...)
}

// Followings returns the list of users the origin is following.
func (c *ConnectionController) Followings(
	app *v04_entity.Application,
	originID uint64,
) (user.List, error) {
	cs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{originID},
		States:  []connection.State{connection.StateConfirmed},
		Types:   []connection.Type{connection.TypeFollow},
	})
	if err != nil {
		return nil, err
	}

	return user.ListFromIDs(c.users, app, cs.ToIDs()...)
}

// Friends returns the list of users the origin is friends with.
func (c *ConnectionController) Friends(
	app *v04_entity.Application,
	originID uint64,
) (user.List, error) {
	fs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{originID},
		States:  []connection.State{connection.StateConfirmed},
		Types:   []connection.Type{connection.TypeFriend},
	})
	if err != nil {
		return nil, err
	}

	ts, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		ToIDs:   []uint64{originID},
		States:  []connection.State{connection.StateConfirmed},
		Types:   []connection.Type{connection.TypeFriend},
	})
	if err != nil {
		return nil, err
	}

	return user.ListFromIDs(c.users, app, append(fs.ToIDs(), ts.FromIDs()...)...)
}

// Update transitions the passed Connection its new state.
func (c *ConnectionController) Update(
	app *v04_entity.Application,
	new *connection.Connection,
) (*connection.Connection, error) {
	_, errs := c.users.Read(app.OrgID, app.ID, new.ToID, false)
	if errs != nil {
		if errs[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
			return nil, ErrNotFound
		}
		return nil, errs[0]
	}

	var (
		fromIDs = []uint64{new.FromID}
		toIDs   = []uint64{new.ToID}
	)

	if new.Type == connection.TypeFriend {
		fromIDs = []uint64{new.FromID, new.ToID}
		toIDs = []uint64{new.FromID, new.ToID}
	}

	cs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: fromIDs,
		ToIDs:   toIDs,
		Types:   []connection.Type{new.Type},
	})
	if err != nil {
		return nil, err
	}

	new.Enabled = true

	if len(cs) > 0 {
		old := cs[0]

		new.FromID = old.FromID
		new.ToID = old.ToID

		if err := validateConTransition(cs[0], new); err != nil {
			return nil, err
		}
	}

	return c.connections.Put(app.Namespace(), new)
}

func validateConTransition(old, new *connection.Connection) error {
	if old.FromID != new.FromID {
		return wrapError(
			ErrInvalidEntity,
			"from id miss-match %d != %d",
			old.FromID,
			new.FromID,
		)
	}

	if old.ToID != new.ToID {
		return wrapError(
			ErrInvalidEntity,
			"to id miss-match %d != %d",
			old.ToID,
			new.ToID,
		)
	}

	if old.Type != new.Type {
		return wrapError(
			ErrInvalidEntity,
			"type miss-match %s != %s",
			string(old.Type),
			string(new.Type),
		)
	}

	switch old.State {
	case connection.StatePending:
		if (new.State == connection.StateConfirmed) ||
			(new.State == connection.StateRejected) {
			return nil
		}
	case connection.StateConfirmed:
		if new.State == connection.StateRejected {
			return nil
		}
	}

	return wrapError(
		ErrInvalidEntity,
		"invalid state transition from %s to %s",
		string(old.State),
		string(new.State),
	)
}

type relation struct {
	isFriend    bool
	isFollower  bool
	isFollowing bool
}

func queryRelation(
	s connection.Service,
	app *v04_entity.Application,
	origin, user uint64,
) (*relation, error) {
	cs, err := s.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			origin,
			user,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		ToIDs: []uint64{
			origin,
			user,
		},
	})
	if err != nil {
		return nil, err
	}

	r := &relation{}

	for _, c := range cs {
		if c.Type == connection.TypeFriend {
			r.isFriend = true
		}

		if c.Type == connection.TypeFollow && c.FromID == origin {
			r.isFollowing = true
		}

		if c.Type == connection.TypeFollow && c.ToID == origin {
			r.isFollower = true
		}
	}

	return r, nil
}
