/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	event struct {
		appUser core.ApplicationUser
		storage core.Event
	}
)

func (evt *event) Read(ctx *context.Context) (err errors.Error) {
	var (
		event   = &entity.Event{}
		eventID string
	)

	eventID = ctx.Vars["eventID"]

	if event, err = evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		eventID); err != nil {
		return
	}

	server.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Update(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
	var (
		eventID string
		er      error
	)

	eventID = ctx.Vars["eventID"]

	existingEvent, err := evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return errors.NewBadRequestError("failed to update the event (2)\n"+er.Error(), er.Error())
	}

	event.ID = eventID
	event.UserID = ctx.Bag["applicationUserID"].(string)

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := evt.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		*existingEvent,
		event,
		true)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
	return
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err errors.Error) {
	var (
		eventID string
		er      error
	)

	eventID = ctx.Vars["eventID"]

	existingEvent, err := evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return errors.NewBadRequestError("failed to update the event (2)\n"+er.Error(), er.Error())
	}

	event.ID = eventID
	event.UserID = ctx.Bag["applicationUserID"].(string)

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := evt.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		*existingEvent,
		event,
		true)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
	return
}

func (evt *event) Delete(ctx *context.Context) (err errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Bag["applicationUserID"].(string)
	eventID := ctx.Vars["eventID"]

	event, err := evt.storage.Read(accountID, applicationID, userID, eventID)
	if err != nil {
		return
	}

	if err = evt.storage.Delete(
		accountID,
		applicationID,
		event); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (evt *event) List(ctx *context.Context) (err errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Bag["applicationUserID"].(string)

	exists, err := evt.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return errors.NewNotFoundError("user not found", "user not found")
	}

	var events []*entity.Event

	if events, err = evt.storage.List(accountID, applicationID, userID); err != nil {
		return
	}

	status := http.StatusOK
	if len(events) == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, events, status, 10)
	return
}

func (evt *event) CurrentUserList(ctx *context.Context) (err errors.Error) {
	var events []*entity.Event

	if events, err = evt.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string)); err != nil {
		return
	}

	status := http.StatusOK
	if len(events) == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, events, status, 10)
	return
}

func (evt *event) Feed(ctx *context.Context) (err errors.Error) {
	var events = []*entity.Event{}

	if events, err = evt.storage.ConnectionList(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string)); err != nil {
		return
	}

	status := http.StatusOK
	if len(events) == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, events, status, 10)
	return
}

func (evt *event) Create(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return errors.NewBadRequestError("failed to create the event (1)\n"+er.Error(), er.Error())
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
		event,
		true); err != nil {
		return
	}

	server.WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err errors.Error) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return errors.NewBadRequestError("failed to create the event (1)\n"+er.Error(), er.Error())
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
		event,
		true); err != nil {
		return
	}

	server.WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

func (evt *event) SearchGeo(ctx *context.Context) (err errors.Error) {
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		er                          error
	)

	if latitude, er = strconv.ParseFloat(ctx.Vars["latitude"], 64); er != nil {
		return errors.NewBadRequestError("failed to read the event by geo (1)\n"+er.Error(), er.Error())
	}

	if longitude, er = strconv.ParseFloat(ctx.Vars["longitude"], 64); er != nil {
		return errors.NewBadRequestError("failed to read the event by geo (2)\n"+er.Error(), er.Error())
	}

	if radius, er = strconv.ParseFloat(ctx.Vars["radius"], 64); er != nil {
		return errors.NewBadRequestError("failed to read the event by geo (3)\n"+er.Error(), er.Error())
	}

	if radius < 1 {
		return errors.NewBadRequestError("failed to read the event by geo (4)\nLocation radius can't be smaller than 2 meters", "radius smaller than 2")
	}

	if events, err = evt.storage.GeoSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), latitude, longitude, radius); err != nil {
		return
	}

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) SearchObject(ctx *context.Context) (err errors.Error) {
	var (
		events    = []*entity.Event{}
		objectKey string
	)

	objectKey = ctx.Vars["objectKey"]

	if events, err = evt.storage.ObjectSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), objectKey); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

func (evt *event) SearchLocation(ctx *context.Context) (err errors.Error) {
	var (
		events   = []*entity.Event{}
		location string
	)

	location = ctx.Vars["location"]

	if events, err = evt.storage.LocationSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), location); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	server.WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

type byIDDesc []*entity.Event

func (a byIDDesc) Len() int           { return len(a) }
func (a byIDDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIDDesc) Less(i, j int) bool { return a[i].ID > a[j].ID }

// NewEventWithApplicationUser returns a new event handler
func NewEventWithApplicationUser(storage core.Event, appUser core.ApplicationUser) server.Event {
	return &event{
		storage: storage,
		appUser: appUser,
	}
}
