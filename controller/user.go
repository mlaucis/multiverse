package controller

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/tapglue/multiverse/platform/generate"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/session"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// UserController bundles the business constraints of Users.
type UserController struct {
	connections connection.Service
	sessions    session.Service
	users       user.Service
}

// NewUserController returns a controller instance.
func NewUserController(
	connections connection.Service,
	sessions session.Service,
	users user.Service,
) *UserController {
	return &UserController{
		connections: connections,
		sessions:    sessions,
		users:       users,
	}
}

// Create stores the provided user and creates a session.
func (c *UserController) Create(
	app *v04_entity.Application,
	origin Origin,
	u *user.User,
) (*user.User, error) {
	if err := constrainUserPrivate(origin, u.Private); err != nil {
		return nil, err
	}

	if err := c.constrainUniqueEmail(app, u); err != nil {
		if !IsInvalidEntity(err) {
			return nil, err
		}

		return c.LoginEmail(app, origin, u.Email, u.Password)
	}

	if err := c.constrainUniqueUsername(app, u); err != nil {
		if !IsInvalidEntity(err) {
			return nil, err
		}

		return c.LoginUsername(app, origin, u.Username, u.Password)
	}

	epw, err := passwordSecure(u.Password)
	if err != nil {
		return nil, err
	}

	u.Enabled = true
	u.Password = epw

	if err := u.Validate(); err != nil {
		return nil, wrapError(ErrInvalidEntity, "%s", err)
	}

	u, err = c.users.Put(app.Namespace(), u)
	if err != nil {
		return nil, err
	}

	err = c.enrichSessionToken(app, u, origin.DeviceID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Delete disables the user.
func (c *UserController) Delete(
	app *v04_entity.Application,
	origin *user.User,
) error {
	origin.Enabled = false
	origin.Deleted = true

	_, err := c.users.Put(app.Namespace(), origin)
	return err
}

// ListByEmails returns all users for the given emails.
func (c *UserController) ListByEmails(
	app *v04_entity.Application,
	originID uint64,
	emails ...string,
) (user.List, error) {
	us, err := c.users.Query(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
		Emails:  emails,
	})
	if err != nil {
		return nil, err
	}

	for _, u := range us {
		r, err := queryRelation(c.connections, app, originID, u.ID)
		if err != nil {
			return nil, err
		}

		u.IsFriend = r.isFriend
		u.IsFollower = r.isFollower
		u.IsFollowing = r.isFollowing
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
	us, err := c.users.Query(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
		SocialIDs: map[string][]string{
			platform: ids,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, u := range us {
		r, err := queryRelation(c.connections, app, originID, u.ID)
		if err != nil {
			return nil, err
		}

		u.IsFriend = r.isFriend
		u.IsFollower = r.isFollower
		u.IsFollowing = r.isFollowing
	}

	return us, nil
}

// LoginEmail finds the user by email and returns it with a valid session token.
func (c *UserController) LoginEmail(
	app *v04_entity.Application,
	origin Origin,
	email string,
	password string,
) (*user.User, error) {
	us, err := c.users.Query(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
		Emails: []string{
			email,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(us) != 1 {
		return nil, ErrNotFound
	}

	return c.login(app, us[0], password, origin.DeviceID)
}

// LoginUsername finds the user by username and returns it with a valid session
// token.
func (c *UserController) LoginUsername(
	app *v04_entity.Application,
	origin Origin,
	username string,
	password string,
) (*user.User, error) {
	us, err := c.users.Query(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
		Usernames: []string{
			username,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(us) != 1 {
		return nil, ErrNotFound
	}

	return c.login(app, us[0], password, origin.DeviceID)
}

// Logout destroys the session stored under token.
func (c *UserController) Logout(
	app *v04_entity.Application,
	origin uint64,
	token string,
) error {
	ss, err := c.sessions.Query(app.Namespace(), session.QueryOptions{
		Enabled: &defaultEnabled,
		IDs: []string{
			token,
		},
		UserIDs: []uint64{
			origin,
		},
	})
	if err != nil {
		return err
	}

	if len(ss) == 0 {
		return nil
	}

	s := ss[0]
	s.Enabled = false

	_, err = c.sessions.Put(app.Namespace(), s)
	return err
}

// Retrieve returns the user with the given id.
func (c *UserController) Retrieve(
	app *v04_entity.Application,
	origin Origin,
	userID uint64,
) (*user.User, error) {
	us, err := c.users.Query(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
		IDs: []uint64{
			userID,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(us) != 1 {
		return nil, ErrNotFound
	}

	u := us[0]

	err = enrichRelation(c.connections, app, origin.UserID, u)
	if err != nil {
		return nil, err
	}

	err = enrichConnectionCounts(c.connections, c.users, app, u)
	if err != nil {
		return nil, err
	}

	if origin.UserID == userID {
		err = c.enrichSessionToken(app, u, origin.DeviceID)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

// Search returns all users for the given query.
func (c *UserController) Search(
	app *v04_entity.Application,
	origin uint64,
	query string,
) (user.List, error) {
	t := []string{query}

	us, err := c.users.Search(app.Namespace(), user.QueryOptions{
		Enabled: &defaultEnabled,
	}, user.SearchOptions{
		Emails:     t,
		Firstnames: t,
		Lastnames:  t,
		Usernames:  t,
	})
	if err != nil {
		return nil, err
	}

	for _, u := range us {
		err = enrichConnectionCounts(c.connections, c.users, app, u)
		if err != nil {
			return nil, err
		}

		err = enrichRelation(c.connections, app, origin, u)
		if err != nil {
			return nil, err
		}
	}

	return us, nil
}

// Update stores the new attributes for the user.
func (c *UserController) Update(
	app *v04_entity.Application,
	origin Origin,
	old *user.User,
	new *user.User,
) (*user.User, error) {
	err := constrainUserPrivate(origin, new.Private)
	if err != nil {
		return nil, err
	}

	new.Enabled = true
	new.ID = old.ID

	if new.Password != "" {
		epw, err := passwordSecure(new.Password)
		if err != nil {
			return nil, err
		}

		new.Password = epw
	} else {
		new.Password = old.Password
	}

	if old.Email != new.Email {
		err := c.constrainUniqueEmail(app, new)
		if err != nil {
			return nil, err
		}
	}

	if new.Private == nil {
		new.Private = old.Private
	}

	if old.Username != new.Username {
		err := c.constrainUniqueUsername(app, new)
		if err != nil {
			return nil, err
		}
	}

	u, err := c.users.Put(app.Namespace(), new)
	if err != nil {
		return nil, err
	}

	err = enrichConnectionCounts(c.connections, c.users, app, u)
	if err != nil {
		return nil, err
	}

	err = c.enrichSessionToken(app, u, origin.DeviceID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *UserController) constrainUniqueEmail(
	app *v04_entity.Application,
	u *user.User,
) error {
	if u.Email != "" {
		us, err := c.users.Query(app.Namespace(), user.QueryOptions{
			Enabled: &defaultEnabled,
			Emails: []string{
				u.Email,
			},
		})
		if err != nil {
			return err
		}

		if len(us) > 0 {
			return wrapError(ErrInvalidEntity, "email in use")
		}
	}

	return nil
}

func (c *UserController) constrainUniqueUsername(
	app *v04_entity.Application,
	u *user.User,
) error {
	if u.Username != "" {
		us, err := c.users.Query(app.Namespace(), user.QueryOptions{
			Enabled: &defaultEnabled,
			Usernames: []string{
				u.Username,
			},
		})
		if err != nil {
			return err
		}

		if len(us) > 0 {
			return wrapError(ErrInvalidEntity, "username in use")
		}
	}

	return nil
}

func enrichConnectionCounts(
	connections connection.Service,
	users user.Service,
	app *v04_entity.Application,
	u *user.User,
) error {
	deleted := false

	cs, err := connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		States: []connection.State{
			connection.StateConfirmed,
		},
		ToIDs: []uint64{
			u.ID,
		},
		Types: []connection.Type{
			connection.TypeFollow,
		},
	})
	if err != nil {
		return err
	}

	if len(cs) > 0 {
		u.FollowerCount, err = users.Count(app.Namespace(), user.QueryOptions{
			Deleted: &deleted,
			Enabled: &defaultEnabled,
			IDs:     cs.FromIDs(),
		})
		if err != nil {
			return err
		}
	}

	cs, err = connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			u.ID,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFollow,
		},
	})
	if err != nil {
		return err
	}

	if len(cs) > 0 {
		u.FollowingCount, err = users.Count(app.Namespace(), user.QueryOptions{
			Deleted: &deleted,
			Enabled: &defaultEnabled,
			IDs:     cs.ToIDs(),
		})
		if err != nil {
			return err
		}
	}

	fs, err := connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			u.ID,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFriend,
		},
	})
	if err != nil {
		return err
	}

	ts, err := connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		States: []connection.State{
			connection.StateConfirmed,
		},
		ToIDs: []uint64{
			u.ID,
		},
		Types: []connection.Type{
			connection.TypeFriend,
		},
	})
	if err != nil {
		return err
	}

	ids := append(fs.ToIDs(), ts.FromIDs()...)

	if len(ids) > 0 {
		u.FriendCount, err = users.Count(app.Namespace(), user.QueryOptions{
			Deleted: &deleted,
			Enabled: &defaultEnabled,
			IDs:     ids,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *UserController) enrichSessionToken(
	app *v04_entity.Application,
	u *user.User,
	deviceID string,
) error {
	ss, err := c.sessions.Query(app.Namespace(), session.QueryOptions{
		DeviceIDs: []string{
			deviceID,
		},
		Enabled: &defaultEnabled,
		UserIDs: []uint64{
			u.ID,
		},
	})
	if err != nil {
		return err
	}

	var s *session.Session

	if len(ss) > 0 {
		s = ss[0]
	} else {
		s, err = c.sessions.Put(app.Namespace(), &session.Session{
			DeviceID: deviceID,
			Enabled:  true,
			UserID:   u.ID,
		})
		if err != nil {
			return err
		}
	}

	u.SessionToken = s.ID

	return nil
}

func (c *UserController) login(
	app *v04_entity.Application,
	u *user.User,
	password string,
	deviceID string,
) (*user.User, error) {
	valid, err := passwordCompare(password, u.Password)
	if err != nil {
		return nil, ErrNotFound
	}

	if !valid {
		return nil, ErrUnauthorized
	}

	err = c.enrichSessionToken(app, u, deviceID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func constrainUserPrivate(origin Origin, private *user.Private) error {
	if !origin.IsBackend() && private != nil {
		return wrapError(
			ErrUnauthorized,
			"private can only be set by backend integration",
		)
	}

	return nil
}

func enrichRelation(
	s connection.Service,
	app *v04_entity.Application,
	origin uint64,
	u *user.User,
) error {
	if origin == u.ID {
		return nil
	}

	r, err := queryRelation(s, app, origin, u.ID)
	if err != nil {
		return err
	}

	u.IsFriend = r.isFriend
	u.IsFollower = r.isFollower
	u.IsFollowing = r.isFollowing

	return nil
}

func passwordCompare(dec, enc string) (bool, error) {
	d, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return false, err
	}

	ps := strings.SplitN(string(d), ":", 3)

	epw, err := base64.StdEncoding.DecodeString(ps[2])
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(ps[0])
	if err != nil {
		return false, err
	}

	ts, err := base64.StdEncoding.DecodeString(ps[1])
	if err != nil {
		return false, err
	}

	esalt := []byte{}
	esalt = append(esalt, []byte(salt)...)
	esalt = append(esalt, []byte(":")...)
	esalt = append(esalt, []byte(ts)...)

	ipw, err := generate.EncryptPassword([]byte(dec), esalt)
	if err != nil {
		return false, err
	}

	return string(epw) == string(ipw), nil
}

func passwordSecure(pw string) (string, error) {
	// create Salt
	salt, err := generate.Salt()
	if err != nil {
		return "", err
	}

	// create scrypt salt
	var (
		esalt = []byte{}
		ts    = []byte(time.Now().Format(time.RFC3339))
	)

	esalt = append(esalt, salt...)
	esalt = append(esalt, []byte(":")...)
	esalt = append(esalt, ts...)

	// encrypt
	epw, err := generate.EncryptPassword([]byte(pw), esalt)
	if err != nil {
		return "", err
	}

	// encode
	enc := fmt.Sprintf(
		"%s:%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(ts),
		base64.StdEncoding.EncodeToString(epw),
	)

	return base64.StdEncoding.EncodeToString([]byte(enc)), nil
}
