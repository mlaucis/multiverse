package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/tgflake"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type (
	event struct {
		appUser core.ApplicationUser
		storage core.Event
	}
)

func (evt *event) CurrentUserRead(ctx *context.Context) (err []errors.Error) {
	var event = &entity.Event{}

	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	if event, err = evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64),
		eventID); err != nil {
		return
	}

	response.ComputeEventLastModified(ctx, event)

	response.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Read(ctx *context.Context) (err []errors.Error) {
	var event = &entity.Event{}

	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	if event, err = evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		userID,
		eventID); err != nil {
		return
	}

	response.ComputeEventLastModified(ctx, event)

	response.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}

	/*	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
		if er != nil {
			return []errors.Error{errmsg.ErrEventIDInvalid}
		}

		userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
		if er != nil {
			return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
		}

		existingEvent, err := evt.storage.Read(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			userID,
			ctx.Bag["applicationUserID"].(uint64),
			eventID)
		if err != nil {
			return
		}

		event := *existingEvent
		if er = json.Unmarshal(ctx.Body, &event); er != nil {
			return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
		}

		event.ID = eventID
		event.UserID = ctx.Bag["applicationUserID"].(uint64)

		if err = validator.UpdateEvent(existingEvent, &event); err != nil {
			return
		}

		updatedEvent, err := evt.storage.Update(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			ctx.Bag["applicationUserID"].(uint64),
			*existingEvent,
			event,
			true)
		if err != nil {
			return
		}

		response.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
		return*/
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err []errors.Error) {
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	existingEvent, err := evt.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64),
		eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	event.ID = eventID
	event.UserID = ctx.Bag["applicationUserID"].(uint64)

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := evt.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64),
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
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	if err = evt.storage.Delete(
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
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Bag["applicationUserID"].(uint64)
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid}
	}

	if err = evt.storage.Delete(
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
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	exists, err := evt.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound}
	}

	var events []*entity.Event

	if events, err = evt.storage.List(accountID, applicationID, userID, userID); err != nil {
		return
	}

	resp := struct {
		Events      []*entity.Event `json:"events"`
		EventsCount int             `json:"events_count"`
	}{
		Events:      events,
		EventsCount: len(events),
	}

	status := http.StatusOK
	response.ComputeEventsLastModified(ctx, resp.Events)

	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	var events []*entity.Event

	if events, err = evt.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64),
		ctx.Bag["applicationUserID"].(uint64)); err != nil {
		return
	}

	resp := struct {
		Events      []*entity.Event `json:"events"`
		EventsCount int             `json:"events_count"`
	}{
		Events:      events,
		EventsCount: len(events),
	}

	status := http.StatusOK
	response.ComputeEventsLastModified(ctx, resp.Events)

	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) Feed(ctx *context.Context) (err []errors.Error) {
	resp := entity.EventsResponseWithUnread{}

	if resp.UnreadCount, resp.Events, err = evt.storage.UserFeed(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	resp.EventsCount = len(resp.Events)

	status := http.StatusOK
	response.ComputeEventsLastModified(ctx, resp.Events)

	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}

	/*	var (
			event = &entity.Event{}
			er    error
		)

		if er = json.Unmarshal(ctx.Body, event); er != nil {
			return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
		}

		event.UserID = ctx.Bag["applicationUserID"].(uint64)
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
			ctx.Bag["applicationUserID"].(uint64),
			event,
			true); err != nil {
			return
		}

		ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/users/events/%d", ctx.Bag["applicationUserID"].(uint64), event.ID))
		response.WriteResponse(ctx, event, http.StatusCreated, 0)
		return*/
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err []errors.Error) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	event.UserID = ctx.Bag["applicationUserID"].(uint64)
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

	event.ID, er = tgflake.FlakeNextID(ctx.Bag["applicationID"].(int64), "events")
	if er != nil {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(er.Error())}
	}

	if event, err = evt.storage.Create(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64),
		event,
		true); err != nil {
		return
	}

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/me/events/%d", event.ID))
	response.WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

func (evt *event) Search(ctx *context.Context) (err []errors.Error) {
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		nearest                     int64
		er                          error
	)

	if l := ctx.Query.Get("lat"); l != "" {
		if latitude, er = strconv.ParseFloat(l, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error())}
		}
	}

	if l := ctx.Query.Get("lon"); l != "" {
		if longitude, er = strconv.ParseFloat(l, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error())}
		}
	}

	if rad := ctx.Query.Get("rad"); rad != "" {
		if radius, er = strconv.ParseFloat(rad, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error())}
		}
	}

	if near := ctx.Query.Get("nearest"); near != "" {
		if nearest, er = strconv.ParseInt(near, 10, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error())}
		}

		if nearest < 1 || nearest > 200 {
			return []errors.Error{errmsg.ErrEventNearestNotInBounds}
		}
	}

	if ctx.Query.Get("lat") != "" && ctx.Query.Get("lon") != "" {
		if radius == 0 && nearest == 0 {
			return []errors.Error{errmsg.ErrEventGeoRadiusAndNearestMissing}
		}

		if radius < 2 && nearest == 0 {
			return []errors.Error{errmsg.ErrEventGeoRadiusUnder2M}
		}

		if radius == 0 && nearest > 200 {
			return []errors.Error{errmsg.ErrEventNearestNotInBounds}
		}

		events, err = evt.storage.GeoSearch(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			ctx.Bag["applicationUserID"].(uint64),
			latitude,
			longitude,
			radius,
			nearest)
	} else if location := ctx.Query.Get("location"); location != "" {
		if events, err = evt.storage.LocationSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(uint64), location); err != nil {
			return
		}
	} else if objectKey := ctx.Query.Get("object"); objectKey != "" {
		if events, err = evt.storage.ObjectSearch(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(uint64), objectKey); err != nil {
			return
		}
	} else {
		err = []errors.Error{errmsg.ErrServerReqNoKnownSearchTermsSupplied}
	}
	if err != nil {
		return
	}

	resp := entity.EventsResponse{
		Events:      events,
		EventsCount: len(events),
	}

	if resp.EventsCount != 0 {
		resp.Users, err = evt.usersFromEvents(ctx, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	status := http.StatusOK
	response.ComputeEventsLastModified(ctx, resp.Events)

	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) UnreadFeed(ctx *context.Context) (err []errors.Error) {
	resp := entity.EventsResponseWithUnread{}

	if resp.UnreadCount, resp.Events, err = evt.storage.UnreadFeed(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	resp.EventsCount = len(resp.Events)

	status := http.StatusOK
	response.ComputeEventsLastModified(ctx, resp.Events)

	if resp.UnreadCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err []errors.Error) {
	count := struct {
		Count int `json:"unread_events_count"`
	}{}

	if count.Count, err = evt.storage.UnreadFeedCount(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	// TODO Maybe not the best idea?
	response.ComputeLastModifiedNow(ctx)

	response.WriteResponse(ctx, count, http.StatusOK, 10)
	return
}

func (evt *event) usersFromEvents(ctx *context.Context, events []*entity.Event) (users map[string]*entity.ApplicationUser, err []errors.Error) {
	users = map[string]*entity.ApplicationUser{}
	eventUsers := map[uint64]bool{}
	for idx := range events {
		eventUsers[events[idx].UserID] = true
	}

	usrs := []uint64{}
	for idx := range eventUsers {
		usrs = append(usrs, idx)
	}

	urs, err := evt.appUser.ReadMultiple(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), usrs)
	if err != nil {
		return
	}

	for idx := range urs {
		users[strconv.FormatUint(urs[idx].ID, 10)] = urs[idx]
	}

	response.SanitizeApplicationUsersMap(users)

	return
}

// NewEvent returns a new event handler
func NewEvent(storage core.Event, appUser core.ApplicationUser) handlers.Event {
	return &event{
		storage: storage,
		appUser: appUser,
	}
}
