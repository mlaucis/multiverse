package controller

import (
	"math/rand"
	"testing"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestEventCreateConstrainVisibility(t *testing.T) {
	var (
		app, owner, c = testSetupEventController(t)
		origin        = Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		}
	)

	_, err := c.Create(app, origin, &event.Event{
		Visibility: event.VisibilityGlobal,
	})

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestEventUpdateConstrainVisibility(t *testing.T) {
	var (
		app, owner, c = testSetupEventController(t)
		origin        = Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		}
	)

	created, err := c.Create(app, origin, testEvent(owner.ID))
	if err != nil {
		t.Fatal(err)
	}

	created.Visibility = event.VisibilityGlobal

	_, err = c.Update(app, origin, created.ID, created)

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testEvent(ownerID uint64) *event.Event {
	return &event.Event{
		UserID:     ownerID,
		Visibility: event.VisibilityConnection,
	}
}

func testSetupEventController(
	t *testing.T,
) (*v04_entity.Application, *user.User, *EventController) {
	var (
		app = &v04_entity.Application{
			ID:    rand.Int63(),
			OrgID: rand.Int63(),
		}
		connections = connection.NewMemService()
		events      = event.NewMemService()
		objects     = object.NewMemService()
		users       = user.NewMemService()
		u           = &user.User{
			ID: uint64(rand.Int63()),
		}
	)

	return app, u, NewEventController(connections, events, objects, users)
}
