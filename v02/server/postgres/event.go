/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"net/http"
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
		event           = &entity.Event{}
		eventID, userID string
	)

	eventID = ctx.Vars["eventID"]
	if !validator.IsValidUUID5(eventID) {
		return invalidEventIDError
	}

	userID = ctx.Vars["applicationUserID"]
	if !validator.IsValidUUID5(userID) {
		return invalidUserIDError
	}

	if event, err = evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		userID,
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
		eventID, userID string
		er              error
	)

	eventID = ctx.Vars["eventID"]
	if !validator.IsValidUUID5(eventID) {
		return invalidEventIDError
	}

	userID = ctx.Vars["applicationUserID"]
	if !validator.IsValidUUID5(userID) {
		return invalidUserIDError
	}

	existingEvent, err := evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		userID,
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

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err errors.Error) {
	var (
		eventID string
		er      error
	)

	eventID = ctx.Vars["eventID"]
	if !validator.IsValidUUID5(eventID) {
		return invalidEventIDError
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

func (evt *event) Delete(ctx *context.Context) (err errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Bag["applicationUserID"].(string)
	eventID := ctx.Vars["eventID"]
	if !validator.IsValidUUID5(eventID) {
		return invalidEventIDError
	}

	event, err := evt.storage.Read(accountID, applicationID, userID, userID, eventID)
	if err != nil {
		return
	}

	if err = evt.storage.Delete(
		accountID,
		applicationID,
		userID,
		event); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (evt *event) List(ctx *context.Context) (err errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Vars["applicationUserID"]
	if !validator.IsValidUUID5(userID) {
		return invalidUserIDError
	}

	exists, err := evt.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return errors.NewNotFoundError("user not found", "user not found")
	}

	var events []*entity.Event

	if events, err = evt.storage.List(accountID, applicationID, userID, userID); err != nil {
		return
	}

	response := struct {
		Events      []*entity.Event `json:"events"`
		EventsCount int             `json:"events_count"`
	}{
		Events:      events,
		EventsCount: len(events),
	}

	status := http.StatusOK
	if response.EventsCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (evt *event) CurrentUserList(ctx *context.Context) (err errors.Error) {
	var events []*entity.Event

	if events, err = evt.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		ctx.Bag["applicationUserID"].(string)); err != nil {
		return
	}

	response := struct {
		Events      []*entity.Event `json:"events"`
		EventsCount int             `json:"events_count"`
	}{
		Events:      events,
		EventsCount: len(events),
	}

	status := http.StatusOK
	if response.EventsCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (evt *event) Feed(ctx *context.Context) (err errors.Error) {
	response := entity.EventsResponseWithUnread{}

	if response.UnreadCount, response.Events, err = evt.storage.UserFeed(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	response.EventsCount = len(response.Events)

	status := http.StatusOK
	if response.EventsCount == 0 {
		status = http.StatusNoContent
	} else {
		response.Users, err = evt.usersFromEvents(ctx, response.Events)
		if err != nil {
			return
		}

		response.UsersCount = len(response.Users)
	}

	server.WriteResponse(ctx, response, status, 10)
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
	if event.Visibility == 0 {
		event.Visibility = entity.EventPublic
	}

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

func (evt *event) CurrentUserCreate(ctx *context.Context) (err errors.Error) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return errors.NewBadRequestError("failed to create the event (1)\n"+er.Error(), er.Error())
	}

	event.UserID = ctx.Bag["applicationUserID"].(string)
	if event.Visibility == 0 {
		event.Visibility = entity.EventPublic
	}

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

func (evt *event) Search(ctx *context.Context) (err errors.Error) {
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		nearest                     int64
		er                          error
	)

	if l := ctx.Query.Get("lat"); l != "" {
		if latitude, er = strconv.ParseFloat(l, 64); er != nil {
			return errors.NewBadRequestError("failed to read the event by geo (1)\n"+er.Error(), er.Error())
		}
	}

	if l := ctx.Query.Get("lon"); l != "" {
		if longitude, er = strconv.ParseFloat(l, 64); er != nil {
			return errors.NewBadRequestError("failed to read the event by geo (2)\nyou must supply a latitude", er.Error())
		}
	}

	if rad := ctx.Query.Get("rad"); rad != "" {
		if radius, er = strconv.ParseFloat(rad, 64); er != nil {
			return errors.NewBadRequestError("failed to read the event by geo (3)\n"+er.Error(), er.Error())
		}
	}

	if near := ctx.Query.Get("nearest"); near != "" {
		if nearest, er = strconv.ParseInt(near, 10, 64); er != nil {
			return errors.NewBadRequestError("failed to read the event by geo (4)\n"+er.Error(), er.Error())
		}

		if nearest < 1 || nearest > 200 {
			return errors.NewBadRequestError("failed to read the events by geo(4)\nnear events limits not within accepted bounds", "nearest not within bounds")
		}
	}

	if ctx.Query.Get("lat") != "" && ctx.Query.Get("lon") != "" {
		if radius == 0 && nearest == 0 {
			return errors.NewBadRequestError("failed to read the event by geo(5) \nyou must specify either a radius or a how many nearest events you want", "invalid radius and nearest")
		}

		if radius < 2 && nearest == 0 {
			return errors.NewBadRequestError("failed to read the event by geo (6)\nLocation radius can't be smaller than 2 meters", "radius smaller than 2")
		}

		if radius == 0 && nearest > 200 {
			return errors.NewBadRequestError("failed to read the event by geo (7)\ncan't have more than 200 nearest events", "nearest is bigger than 200")
		}

		events, err = evt.storage.GeoSearch(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			ctx.Bag["applicationUserID"].(string),
			latitude,
			longitude,
			radius,
			nearest)
	} else if location := ctx.Query.Get("location"); location != "" {
		if events, err = evt.storage.LocationSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(string), location); err != nil {
			return
		}
	} else if objectKey := ctx.Query.Get("object"); objectKey != "" {
		if events, err = evt.storage.ObjectSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(string), objectKey); err != nil {
			return
		}
	} else {
		err = errors.NewBadRequestError("failed to search for events\nno known search terms supplied", "failed to search for events\nno known search terms supplied")
	}
	if err != nil {
		return
	}

	response := entity.EventsResponse{
		Events:      events,
		EventsCount: len(events),
	}

	if response.EventsCount != 0 {
		response.Users, err = evt.usersFromEvents(ctx, response.Events)
		if err != nil {
			return
		}

		response.UsersCount = len(response.Users)
	}

	status := http.StatusOK
	if response.EventsCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (evt *event) UnreadFeed(ctx *context.Context) (err errors.Error) {
	response := entity.EventsResponseWithUnread{}

	if response.UnreadCount, response.Events, err = evt.storage.UnreadFeed(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	response.EventsCount = len(response.Events)

	status := http.StatusOK
	if response.UnreadCount == 0 {
		status = http.StatusNoContent
	} else {
		response.Users, err = evt.usersFromEvents(ctx, response.Events)
		if err != nil {
			return
		}

		response.UsersCount = len(response.Users)
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err errors.Error) {
	count := struct {
		Count int `json:"unread_events_count"`
	}{}

	if count.Count, err = evt.storage.UnreadFeedCount(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, count, http.StatusOK, 10)
	return
}

func (evt *event) usersFromEvents(ctx *context.Context, events []*entity.Event) (users map[string]*entity.ApplicationUser, err errors.Error) {
	users = map[string]*entity.ApplicationUser{}
	eventUsers := map[string]bool{}
	for idx := range events {
		eventUsers[events[idx].UserID] = true
	}

	usrs := []string{}
	for idx := range eventUsers {
		usrs = append(usrs, idx)
	}

	urs, err := evt.appUser.ReadMultiple(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), usrs)
	if err != nil {
		return
	}

	for idx := range urs {
		users[urs[idx].ID] = urs[idx]
	}

	sanitizeApplicationUsersMap(users)

	return
}

// NewEventWithApplicationUser returns a new event handler
func NewEventWithApplicationUser(storage core.Event, appUser core.ApplicationUser) server.Event {
	return &event{
		storage: storage,
		appUser: appUser,
	}
}
