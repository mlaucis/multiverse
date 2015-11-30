package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/response"
	"github.com/tapglue/multiverse/v04/validator"
)

type event struct {
	appUser core.ApplicationUser
	conn    core.Connection
	storage core.Event
}

func (evt *event) CurrentUserRead(ctx *context.Context) (err []errors.Error) {
	var event = &entity.Event{}

	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	if event, err = evt.storage.Read(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		eventID); err != nil {
		return
	}

	response.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Read(ctx *context.Context) (err []errors.Error) {
	var event = &entity.Event{}

	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	if event, err = evt.storage.Read(
		ctx.OrganizationID,
		ctx.ApplicationID,
		userID,
		eventID); err != nil {
		return
	}

	response.WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

func (evt *event) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}

	/*	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
		if er != nil {
			return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
		}

		userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
		if er != nil {
			return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
		}

		existingEvent, err := evt.storage.Read(
			ctx.OrganizationID,
			ctx.ApplicationID,
			userID,
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

		updatedEvent, err := evt.storage.Update(
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
		return*/
}

func (evt *event) CurrentUserUpdate(ctx *context.Context) (err []errors.Error) {
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
	}

	existingEvent, err := evt.storage.Read(
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

	updatedEvent, err := evt.storage.Update(
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
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID := ctx.ApplicationUserID
	eventID, er := strconv.ParseUint(ctx.Vars["eventID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrEventIDInvalid.SetCurrentLocation()}
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
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := evt.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
	}

	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}
	var events []*entity.Event
	if events, err = evt.storage.ListUser(accountID, applicationID, userID, ctx.ApplicationUserID, conditions); err != nil {
		return
	}

	resp := entity.EventsResponse{
		Events:      evt.presentationEvent(events),
		EventsCount: len(events),
	}

	status := http.StatusOK
	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, userID, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}

	events := []*entity.Event{}
	if events, err = evt.storage.ListUser(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		ctx.ApplicationUserID,
		conditions); err != nil {
		return
	}

	resp := entity.EventsResponse{
		Events:      evt.presentationEvent(events),
		EventsCount: len(events),
	}

	status := http.StatusOK
	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, ctx.ApplicationUserID, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) Feed(ctx *context.Context) (err []errors.Error) {
	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}

	resp := entity.EventsResponseWithUnread{}
	events := []*entity.Event{}
	if resp.UnreadCount, events, err = evt.storage.UserFeed(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser,
		conditions); err != nil {
		return
	}

	resp.Events = evt.presentationEvent(events)
	resp.EventsCount = len(resp.Events)

	status := http.StatusOK

	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, ctx.ApplicationUserID, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}

	/*	var (
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
			evt.appUser,
			ctx.OrganizationID,
			ctx.ApplicationID,
			event); err != nil {
			return
		}

		if event, err = evt.storage.Create(
			ctx.OrganizationID,
			ctx.ApplicationID,
			ctx.ApplicationUserID,
			event,
			true); err != nil {
			return
		}

		ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/users/events/%d", ctx.ApplicationUserID, event.ID))
		response.WriteResponse(ctx, event, http.StatusCreated, 0)
		return*/
}

func (evt *event) CurrentUserCreate(ctx *context.Context) (err []errors.Error) {
	var (
		pe = &entity.PresentationEvent{}
		er error
	)

	if er = json.Unmarshal(ctx.Body, pe); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	ev := pe.Event

	ev.UserID = ctx.ApplicationUserID
	if pe.Visibility == 0 {
		ev.Visibility = entity.EventPublic
	}

	if err = validator.CreateEvent(
		evt.appUser,
		ctx.OrganizationID,
		ctx.ApplicationID,
		ev); err != nil {
		return
	}

	ev.ID, er = tgflake.FlakeNextID(ctx.ApplicationID, "events")
	if er != nil {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(er.Error()).SetCurrentLocation()}
	}

	err = evt.storage.Create(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		ev)
	if err != nil {
		return
	}

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/me/events/%d", ev.ID))
	response.WriteResponse(ctx, &entity.PresentationEvent{Event: pe.Event}, http.StatusCreated, 0)
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
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error()).SetCurrentLocation()}
		}
	}

	if l := ctx.Query.Get("lon"); l != "" {
		if longitude, er = strconv.ParseFloat(l, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error()).SetCurrentLocation()}
		}
	}

	if rad := ctx.Query.Get("rad"); rad != "" {
		if radius, er = strconv.ParseFloat(rad, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error()).SetCurrentLocation()}
		}
	}

	if near := ctx.Query.Get("nearest"); near != "" {
		if nearest, er = strconv.ParseInt(near, 10, 64); er != nil {
			return []errors.Error{errmsg.ErrServerReqParseFloat.UpdateMessage(er.Error()).SetCurrentLocation()}
		}

		if nearest < 1 || nearest > 200 {
			return []errors.Error{errmsg.ErrEventNearestNotInBounds.SetCurrentLocation()}
		}
	}

	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}

	if conditions != nil {
		events, err = evt.storage.List(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID, conditions)
	} else if ctx.Query.Get("lat") != "" && ctx.Query.Get("lon") != "" {
		if radius == 0 && nearest == 0 {
			return []errors.Error{errmsg.ErrEventGeoRadiusAndNearestMissing.SetCurrentLocation()}
		}

		if radius < 2 && nearest == 0 {
			return []errors.Error{errmsg.ErrEventGeoRadiusUnder2M.SetCurrentLocation()}
		}

		if radius == 0 && nearest > 200 {
			return []errors.Error{errmsg.ErrEventNearestNotInBounds.SetCurrentLocation()}
		}

		events, err = evt.storage.GeoSearch(
			ctx.OrganizationID,
			ctx.ApplicationID,
			ctx.ApplicationUserID,
			latitude,
			longitude,
			radius,
			nearest)
	} else if location := ctx.Query.Get("location"); location != "" {
		events, err = evt.storage.LocationSearch(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID, location)
	} else {
		err = []errors.Error{errmsg.ErrServerReqNoKnownSearchTermsSupplied.SetCurrentLocation()}
	}
	if err != nil {
		return
	}

	resp := entity.EventsResponse{
		Events:      evt.presentationEvent(events),
		EventsCount: len(events),
	}

	if resp.EventsCount != 0 {
		resp.Users, err = evt.usersFromEvents(ctx, ctx.ApplicationUserID, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	status := http.StatusOK
	if resp.EventsCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) UnreadFeed(ctx *context.Context) (err []errors.Error) {
	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}

	resp := entity.EventsResponseWithUnread{}
	events := []*entity.Event{}
	if resp.UnreadCount, events, err = evt.storage.UnreadFeed(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser,
		conditions); err != nil {
		return
	}

	resp.Events = evt.presentationEvent(events)
	resp.EventsCount = len(resp.Events)

	status := http.StatusOK

	if resp.UnreadCount == 0 {
		status = http.StatusNoContent
	} else {
		resp.Users, err = evt.usersFromEvents(ctx, ctx.ApplicationUserID, resp.Events)
		if err != nil {
			return
		}

		resp.UsersCount = len(resp.Users)
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (evt *event) UnreadFeedCount(ctx *context.Context) (err []errors.Error) {
	conditions, err := core.NewEventFilter(ctx.Query.Get("where"))
	if err != nil {
		return err
	}

	count := struct {
		Count int `json:"unread_events_count"`
	}{}
	if count.Count, err = evt.storage.UnreadFeedCount(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser,
		conditions); err != nil {
		return
	}

	response.WriteResponse(ctx, count, http.StatusOK, 10)
	return
}

func (evt *event) usersFromEvents(ctx *context.Context, currentUserID uint64, events []*entity.PresentationEvent) (users map[string]*entity.PresentationApplicationUser, err []errors.Error) {
	users = map[string]*entity.PresentationApplicationUser{}
	eventUsers := map[uint64]bool{}
	for idx := range events {
		eventUsers[events[idx].UserID] = true
		if events[idx].Target != nil &&
			events[idx].Target.Type == "tg_user" {
			if userID, ok := events[idx].Target.ID.(uint64); ok {
				eventUsers[userID] = true
			} else if userID, ok := events[idx].Target.ID.(string); ok {
				if userID, err := strconv.ParseUint(userID, 10, 64); err == nil {
					eventUsers[userID] = true
				}
			}
		}
	}

	usrs := []uint64{}
	for idx := range eventUsers {
		usrs = append(usrs, idx)
	}

	urs, err := evt.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, usrs)
	if err != nil {
		return
	}
	response.SanitizeApplicationUsers(urs)

	for idx := range urs {
		relation, err := evt.conn.Relation(ctx.OrganizationID, ctx.ApplicationID, currentUserID, urs[idx].ID)
		if err != nil {
			return nil, err
		} else if relation != nil {
			urs[idx].Relation = *relation
		}
		users[strconv.FormatUint(urs[idx].ID, 10)] = &entity.PresentationApplicationUser{
			ApplicationUser: urs[idx],
		}

	}

	return
}

func (evt *event) presentationEvent(events []*entity.Event) []*entity.PresentationEvent {
	eventsWithIDString := make([]*entity.PresentationEvent, len(events))
	for idx := range events {
		eventsWithIDString[idx] = &entity.PresentationEvent{
			Event: events[idx],
		}
	}

	return eventsWithIDString
}

// NewEvent returns a new event handler
func NewEvent(storage core.Event, appUser core.ApplicationUser, conn core.Connection) handlers.Event {
	return &event{
		storage: storage,
		appUser: appUser,
		conn:    conn,
	}
}
