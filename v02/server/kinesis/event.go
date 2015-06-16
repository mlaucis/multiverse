/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	event struct {
		appUser core.ApplicationUser
		storage core.Event
	}
)

func (evt *event) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err []errors.Error) {
	var (
		eventID string
		er      error
	)

	eventID = ctx.Vars["eventID"]
	if !validator.IsValidUUID5(eventID) {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	existingEvent, err := evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		ctx.Bag["applicationUserID"].(string),
		eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	event.ID = eventID
	event.UserID = ctx.Bag["applicationUserID"].(string)

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := evt.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		*existingEvent,
		event,
		true)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
	return
}

func (evt *event) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) Feed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) Search(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) SearchGeo(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) SearchObject(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) SearchLocation(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) UnreadFeed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
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
