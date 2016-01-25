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

// ObjectCreate stores the object from the payload in the object.Service.
func ObjectCreate(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
			ro   = &responseObject{}
		)

		err := json.NewDecoder(r.Body).Decode(ro)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		object, err := c.Create(app, ro.Object, user)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &responseObject{Object: object})
	}
}

// ObjectDelete flags the Object as deleted.
func ObjectDelete(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		id, err := strconv.ParseUint(mux.Vars(r)["objectID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, id)
		if err != nil {
			respondError(w, 0, err)
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
			respondError(w, 0, err)
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
			respondError(w, 0, err)
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
			respondError(w, 0, err)
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
			respondError(w, 0, err)
			return
		}

		object, err := c.Retrieve(app, id)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &responseObject{Object: object})
	}
}

// ObjectUpdate replaces an object with new values.
func ObjectUpdate(c *controller.ObjectController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app = appFromContext(ctx)
			ro  = &responseObject{}
		)

		id, err := strconv.ParseUint(mux.Vars(r)["objectID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(ro)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		updated, err := c.Update(app, id, ro.Object)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &responseObject{Object: updated})
	}
}

type responseObject struct {
	*object.Object
}

type responseObjectFields struct {
	Attachments []object.Attachment `json:"attachments"`
	CreatedAt   time.Time           `json:"created_at,omitempty"`
	ID          string              `json:"id"`
	Latitude    float64             `json:"latitude"`
	Location    string              `json:"location"`
	Longitude   float64             `json:"longitude"`
	ObjectID    string              `json:"tg_object_id,omitempty"`
	Tags        []string            `json:"tags"`
	Type        string              `json:"type"`
	UpdatedAt   time.Time           `json:"updated_at,omitempty"`
	UserID      string              `json:"user_id"`
	Visibility  object.Visibility   `json:"visibility"`
}

func (ro *responseObject) MarshalJSON() ([]byte, error) {
	var (
		o = ro.Object
		r = responseObjectFields{
			Attachments: o.Attachments,
			CreatedAt:   o.CreatedAt,
			ID:          strconv.FormatUint(o.ID, 10),
			Latitude:    o.Latitude,
			Location:    o.Location,
			Longitude:   o.Longitude,
			Tags:        o.Tags,
			Type:        o.Type,
			UpdatedAt:   o.UpdatedAt,
			UserID:      strconv.FormatUint(o.OwnerID, 10),
			Visibility:  o.Visibility,
		}
	)

	if o.ObjectID > 0 {
		r.ObjectID = strconv.FormatUint(o.ObjectID, 10)
	}

	return json.Marshal(r)
}

func (ro *responseObject) UnmarshalJSON(raw []byte) error {
	f := responseObjectFields{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	ro.Object = &object.Object{
		Attachments: f.Attachments,
		Latitude:    f.Latitude,
		Location:    f.Location,
		Longitude:   f.Longitude,
		Tags:        f.Tags,
		Type:        f.Type,
		Visibility:  f.Visibility,
	}

	if f.ObjectID != "" {
		oID, err := strconv.ParseUint(f.ObjectID, 10, 64)
		if err != nil {
			return err
		}

		ro.Object.ObjectID = oID
	}

	return nil
}

func fromObjects(os []*object.Object) []*responseObject {
	rs := []*responseObject{}

	for _, o := range os {
		rs = append(rs, &responseObject{Object: o})
	}

	return rs
}
