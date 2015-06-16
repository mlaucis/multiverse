/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
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
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		event   = &entity.Event{}
		eventID string
	)

	eventID = ctx.Vars["eventId"]

	if event, err = evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		ctx.Bag["applicationUserID"].(string),
		eventID); err != nil {
		return
	}

	server.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		eventID string
		er      error
	)

	eventID = ctx.Vars["eventId"]

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
		return []errors.Error{errors.NewBadRequestError(0, "failed to update the event (2)\n"+er.Error(), er.Error())}
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

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

func (evt *event) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	event := &entity.Event{}
	if er := json.Unmarshal(ctx.Body, event); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to delete the event (1)\n"+er.Error(), er.Error())}
	}

	if err = evt.storage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		event); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (evt *event) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var events []*entity.Event

	if events, err = evt.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		ctx.Bag["applicationUserID"].(string)); err != nil {
		return
	}

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

func (evt *event) Feed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var events = []*entity.Event{}

	if _, events, err = evt.storage.UserFeed(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to create the event (1)\n"+er.Error(), er.Error())}
	}

	event.UserID = ctx.Bag["applicationUserID"].(string)

	if err = validator.CreateEvent(
		evt.appUser,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		event); err != nil {
		return
	}

	if event, err = evt.storage.Create(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		event,
		true); err != nil {
		return
	}

	server.WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

func (evt *event) Search(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
}

func (evt *event) SearchGeo(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		er                          error
	)

	if latitude, er = strconv.ParseFloat(ctx.Vars["latitude"], 64); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to read the event by geo (1)\n"+er.Error(), er.Error())}
	}

	if longitude, er = strconv.ParseFloat(ctx.Vars["longitude"], 64); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to read the event by geo (2)\n"+er.Error(), er.Error())}
	}

	if radius, er = strconv.ParseFloat(ctx.Vars["radius"], 64); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to read the event by geo (3)\n"+er.Error(), er.Error())}
	}

	if radius < 2 {
		return []errors.Error{errors.NewBadRequestError(0, "failed to read the event by geo (4)\nLocation radius can't be smaller than 2 meters", "radius smaller than 2")}
	}

	events, err = evt.storage.GeoSearch(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		latitude,
		longitude,
		radius,
		0)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) SearchObject(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		events    = []*entity.Event{}
		objectKey string
	)

	objectKey = ctx.Vars["objectKey"]

	if events, err = evt.storage.ObjectSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(string), objectKey); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) SearchLocation(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		events   = []*entity.Event{}
		location string
	)

	location = ctx.Vars["location"]

	if events, err = evt.storage.LocationSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(string), location); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) UnreadFeed(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

type byIDDesc []*entity.Event

func (a byIDDesc) Len() int           { return len(a) }
func (a byIDDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIDDesc) Less(i, j int) bool { return a[i].ID > a[j].ID }

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
