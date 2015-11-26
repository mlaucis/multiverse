package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/object"
)

// Handler is the gateway specific http.HandlerFunc expecting a context.Context.
type Handler func(context.Context, http.ResponseWriter, *http.Request)

// ObjectCreate stores the object from the payload in the object.Service.
func ObjectCreate(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app     = appFromContext(ctx)
			user    = userFromContext(ctx)
			rObject = &responseObject{}
		)

		err := json.NewDecoder(r.Body).Decode(rObject)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		object, err := toObject(rObject)
		if err != nil {
			respondError(w, http.StatusBadRequest, 0, err)
			return
		}

		object, err = c.Create(app, object, user)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, fromObject(object))
	}
}

// ObjectDelete flags the Object as deleted.
func ObjectDelete(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		id, err := strconv.ParseUint(mux.Vars(r)["objectID"], 10, 64)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		err = c.Delete(app, id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// ObjectList returns all objects for a user.
func ObjectList(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		os, err := c.List(app, user.ID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, struct {
			Objects []*responseObject `json:"objects"`
		}{
			Objects: fromObjects(os),
		})
	}
}

// ObjectListAll returns all objects for the App.
func ObjectListAll(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		os, err := c.ListAll(app)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, struct {
			Objects []*responseObject `json:"objects"`
		}{
			Objects: fromObjects(os),
		})
	}
}

// ObjectListConnections returns objects for all connections of the current
// user.
func ObjectListConnections(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		os, err := c.ListConnections(app, user.ID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, struct {
			Objects []*responseObject `json:"objects"`
		}{
			Objects: fromObjects(os),
		})
	}
}

// ObjectRetrieve returns the requested Object.
func ObjectRetrieve(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		id, err := strconv.ParseUint(mux.Vars(r)["objectID"], 10, 64)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		object, err := c.Retrieve(app, id)
		if err != nil {
			if err == controller.ErrNotFound {
				respondError(w, http.StatusNotFound, 0, err)
				return
			}

			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, fromObject(object))
	}
}

// ObjectUpdate replaces an object with new values.
func ObjectUpdate(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app     = appFromContext(ctx)
			rObject = &responseObject{}
		)

		id, err := strconv.ParseUint(mux.Vars(r)["objectID"], 10, 64)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(rObject)
		if err != nil {
			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		object, err := toObject(rObject)
		if err != nil {
			respondError(w, http.StatusBadRequest, 0, err)
			return
		}

		updated, err := c.Update(app, id, object)
		if err != nil {
			if err == controller.ErrNotFound {
				respondError(w, http.StatusNotFound, 0, err)
				return
			}

			respondError(w, http.StatusInternalServerError, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, fromObject(updated))
	}
}

type responseObject struct {
	Attachments []object.Attachment `json:"attachments"`
	CreatedAt   time.Time           `json:"created_at"`
	ID          string              `json:"id"`
	Latitude    float64             `json:"latitude"`
	Location    string              `json:"location"`
	Longitude   float64             `json:"longitude"`
	ObjectID    string              `json:"tg_object_id,omitempty"`
	Tags        []string            `json:"tags"`
	TargetID    string              `json:"target_id,omitempty"`
	Type        string              `json:"type"`
	UpdatedAt   time.Time           `json:"updated_at"`
	UserID      string              `json:"user_id"`
	Visibility  object.Visibility   `json:"visibility"`
}

func fromObject(o *object.Object) *responseObject {
	r := &responseObject{
		Attachments: o.Attachments,
		CreatedAt:   o.CreatedAt,
		ID:          strconv.FormatUint(o.ID, 10),
		Latitude:    o.Latitude,
		Location:    o.Location,
		Longitude:   o.Longitude,
		Tags:        o.Tags,
		TargetID:    o.TargetID,
		Type:        o.Type,
		UpdatedAt:   o.UpdatedAt,
		UserID:      strconv.FormatUint(o.OwnerID, 10),
		Visibility:  o.Visibility,
	}

	if o.ObjectID > 0 {
		r.ObjectID = strconv.FormatUint(o.ObjectID, 10)
	}

	return r
}

func fromObjects(os []*object.Object) []*responseObject {
	rs := []*responseObject{}

	for _, o := range os {
		rs = append(rs, fromObject(o))
	}

	return rs
}

func toObject(r *responseObject) (*object.Object, error) {
	o := &object.Object{
		Attachments: r.Attachments,
		Latitude:    r.Latitude,
		Location:    r.Location,
		Longitude:   r.Longitude,
		Tags:        r.Tags,
		TargetID:    r.TargetID,
		Type:        r.Type,
		Visibility:  r.Visibility,
	}

	if r.ObjectID != "" {
		oID, err := strconv.ParseUint(r.ObjectID, 10, 64)
		if err != nil {
			return nil, err
		}

		o.ObjectID = oID
	}

	return o, nil
}
