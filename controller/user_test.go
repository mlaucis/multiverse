package controller

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/tapglue/multiverse/platform/generate"
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/session"
	"github.com/tapglue/multiverse/service/user"
)

func TestUserCreateConstrainPrivate(t *testing.T) {
	var (
		app, c = testSetupUserController(t)
		origin = Origin{Integration: IntegrationApplication}
	)

	u := testUser()
	u.Private = &user.Private{
		Verified: true,
	}

	_, err := c.Create(app, origin, u)

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestUserUpdateConstrainPrivate(t *testing.T) {
	var (
		app, c = testSetupUserController(t)
		u      = testUser()
	)

	created, err := c.users.Put(app.Namespace(), u)
	if err != nil {
		t.Fatal(err)
	}

	created.Private = &user.Private{
		Type:     "brand",
		Verified: true,
	}

	_, err = c.Update(
		app,
		Origin{
			Integration: IntegrationApplication,
			UserID:      created.ID,
		},
		u,
		created,
	)

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPassword(t *testing.T) {
	password := "foobar"

	epw, err := passwordSecure(password)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := passwordCompare(password, epw)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := valid, true; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testSetupUserController(
	t *testing.T,
) (*app.App, *UserController) {
	var (
		a = &app.App{
			ID:    uint64(rand.Int63()),
			OrgID: uint64(rand.Int63()),
		}
		connections = connection.NewMemService()
		sessions    = session.NewMemService()
		users       = user.NewMemService()
	)

	return a, NewUserController(connections, sessions, users)
}

func testUser() *user.User {
	return &user.User{
		Email: fmt.Sprintf(
			"user%d@tapglue.test", rand.Int63(),
		),
		Enabled:  true,
		Password: generate.RandomString(8),
	}
}
