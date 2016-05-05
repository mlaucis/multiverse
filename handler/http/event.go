package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
)

// EventCreate stores a new event for the current user.
func EventCreate(c *controller.EventController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadEvent{}
			tokenType   = tokenTypeFromContext(ctx)
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		event, err := c.Create(
			currentApp,
			createOrigin(tokenType, currentUser.ID),
			p.event,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadEvent{event: event})
	}
}

// EventDelete marks an event as disabled.
func EventDelete(c *controller.EventController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, currentUser.ID, id)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// EventListMe returns all events for the current user.
func EventListMe(c *controller.EventController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		opts, err := whereToOpts(r.URL.Query().Get("where"))
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		l, err := c.List(app, currentUser.ID, currentUser.ID, opts)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(l.Events) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedEvents{
			events:  l.Events,
			postMap: l.PostMap,
			userMap: l.UserMap,
		})
	}
}

// EventListUser returns all events as visible by the current user,
func EventListUser(c *controller.EventController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		userID, err := strconv.ParseUint(mux.Vars(r)["userID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		l, err := c.List(app, currentUser.ID, userID, nil)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(l.Events) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedEvents{
			events:  l.Events,
			postMap: l.PostMap,
			userMap: l.UserMap,
		})
	}
}

// EventUpdate replaces an event with new values.
func EventUpdate(c *controller.EventController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadEvent{}
			tokenType   = tokenTypeFromContext(ctx)
		)

		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		event, err := c.Update(
			currentApp,
			createOrigin(tokenType, currentUser.ID),
			id,
			p.event,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadEvent{event: event})
	}
}

type payloadEvent struct {
	event *event.Event
}

func (p *payloadEvent) MarshalJSON() ([]byte, error) {
	f := struct {
		ID           uint64        `json:"id"`
		IDString     string        `json:"id_string"`
		Language     string        `json:"language"`
		Object       *event.Object `json:"object"`
		ObjectID     string        `json:"tg_object_id"`
		Owned        bool          `json:"owned"`
		Target       *event.Target `json:"target,omitempty"`
		Type         string        `json:"type"`
		UserID       uint64        `json:"user_id"`
		UserIDString string        `json:"user_id_string"`
		Visibility   uint8         `json:"visibility"`
		CreatedAt    time.Time     `json:"created_at"`
		UpdatedAt    time.Time     `json:"updated_at"`
	}{
		ID:           p.event.ID,
		IDString:     strconv.FormatUint(p.event.ID, 10),
		Language:     p.event.Language,
		ObjectID:     strconv.FormatUint(p.event.ObjectID, 10),
		Owned:        p.event.Owned,
		Type:         p.event.Type,
		UserID:       p.event.UserID,
		UserIDString: strconv.FormatUint(p.event.UserID, 10),
		Visibility:   uint8(p.event.Visibility),
		CreatedAt:    p.event.CreatedAt,
		UpdatedAt:    p.event.UpdatedAt,
	}

	if p.event.Object != nil {
		f.Object = p.event.Object
	}

	if p.event.Target != nil {
		f.Target = p.event.Target
	}

	return json.Marshal(f)
}

func (p *payloadEvent) UnmarshalJSON(raw []byte) error {
	f := struct {
		Language   string         `json:"language"`
		Object     *payloadObject `json:"object"`
		Target     *event.Target  `json:"target,omitempty"`
		Type       string         `json:"type"`
		Visibility uint8          `json:"visibility"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	e := &event.Event{
		Language:   f.Language,
		Target:     f.Target,
		Type:       f.Type,
		Visibility: event.Visibility(f.Visibility),
	}

	if f.Object != nil {
		e.Object = &event.Object{
			ID:   f.Object.ID,
			Type: f.Object.Type,
			URL:  f.Object.URL,
		}
	}

	p.event = e

	return nil
}

type payloadObject struct {
	ID   string
	Type string
	URL  string
}

func (p *payloadObject) UnmarshalJSON(raw []byte) error {
	f := struct {
		ID   interface{} `json:"id"`
		Type string      `json:"type"`
		URL  string      `json:"url"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.Type = f.Type
	p.URL = f.URL

	id, err := parseID(f.ID)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

func parseID(input interface{}) (string, error) {
	var id string

	switch t := input.(type) {
	case float64:
		id = fmt.Sprintf("%d", int64(t))
	case int:
		id = strconv.Itoa(t)
	case int64:
		id = strconv.FormatInt(t, 10)
	case uint64:
		id = strconv.FormatUint(t, 10)
	case string:
		id = t
	default:
		return "", fmt.Errorf("unexpected value for id")
	}

	return id, nil
}
