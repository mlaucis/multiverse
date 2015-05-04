/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/server"
)

type (
	event struct {
		appUser core.ApplicationUser
		storage core.Event
	}
)

func (evt *event) Read(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) Update(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) Delete(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) List(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) CurrentUserList(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) Feed(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) Create(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) SearchGeo(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) SearchObject(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (evt *event) SearchLocation(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewEvent returns a new event handler
func NewEvent(storage core.Event) server.Event {
	return &event{
		storage: storage,
	}
}

// NewEventWithApplicationUser returns a new event handler
func NewEventWithApplicationUser(storage core.Event, appUser core.ApplicationUser) server.Event {
	return &event{
		storage: storage,
		appUser: appUser,
	}
}
