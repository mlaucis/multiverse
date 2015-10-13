package kinesis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
	"github.com/tapglue/multiverse/v03/validator"
)

type event struct {
	readAppUser  core.ApplicationUser
	writeStorage core.Event
	readStorage  core.Event
}

func (evt *event) CurrentUserRead(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err []errors.Error) {
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	existingEvent, err := evt.readStorage.Read(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	event.ID = eventID
	event.UserID = ctx.ApplicationUserID

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := evt.writeStorage.Update(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		*existingEvent,
		event,
		true)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
	return
}

func (evt *event) Delete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	if err = evt.writeStorage.Delete(
		accountID,
		applicationID,
		userID,
		eventID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (evt *event) CurrentUserDelete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID := ctx.ApplicationUserID
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	if err = evt.writeStorage.Delete(
		accountID,
		applicationID,
		userID,
		eventID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (evt *event) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) Feed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err []errors.Error) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	event.UserID = ctx.ApplicationUserID
	if event.Visibility == 0 {
		event.Visibility = entity.EventPublic
	}

	if err = validator.CreateEvent(
		evt.readAppUser,
		ctx.OrganizationID,
		ctx.ApplicationID,
		event); err != nil {
		return
	}

	event.ID, er = tgflake.FlakeNextID(ctx.ApplicationID, "events")
	if er != nil {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(er.Error()).SetCurrentLocation()}
	}

	ev, err := evt.writeStorage.Create(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		event,
		true)
	if err != nil {
		return
	}

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/me/events/%d", ev.ID))
	response.WriteResponse(ctx, ev, http.StatusCreated, 0)
	return
}

func (evt *event) Search(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) SearchGeo(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) SearchObject(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) SearchLocation(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) UnreadFeed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// NewEventWithApplicationUser returns a new event handler
func NewEventWithApplicationUser(witeStorage, readStorage core.Event, readAppUser core.ApplicationUser) handlers.Event {
	return &event{
		writeStorage: witeStorage,
		readStorage:  readStorage,
		readAppUser:  readAppUser,
	}
}
