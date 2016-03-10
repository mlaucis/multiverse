package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"

	"golang.org/x/net/context"
)

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

		l, err := c.ListUser(app, currentUser.ID, userID)
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
