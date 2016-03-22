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
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/user"
)

// ConnectionByState returns all connections for a user for a certain state.
func ConnectionByState(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		f, err := c.ByState(app, currentUser.ID, connection.State(mux.Vars(r)["state"]))
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(f.Connections) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadConnections{
			cons:    f.Connections,
			origin:  currentUser.ID,
			userMap: f.UserMap,
		})
	}
}

// ConnectionDelete flags the given connection as disabled.
func ConnectionDelete(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		toID, err := strconv.ParseUint(mux.Vars(r)["toID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		con := &connection.Connection{
			FromID: currentUser.ID,
			ToID:   toID,
			State:  connection.StatePending,
			Type:   connection.Type(mux.Vars(r)["type"]),
		}

		if err := con.Validate(); err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, con)
		if err != nil {
			if controller.IsInvalidEntity(err) {
				respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			} else {
				respondError(w, 0, err)
			}

			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// ConnectionFollowers returns the list of users who follow the current user.
func ConnectionFollowers(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.Followers(app, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// ConnectionFollowings returns the list of users the current user is following.
func ConnectionFollowings(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.Followings(app, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// ConnectionFriends returns the list of users the current user is friends with.
func ConnectionFriends(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.Friends(app, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// ConnectionSocial takes a list of connection ids and creates connections for
// the given user.
func ConnectionSocial(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadSocial{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		us, err := c.CreateSocial(
			app,
			currentUser.ID,
			p.Type,
			p.State,
			p.Platform,
			p.ConnectionIDs...,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// ConnectionUpdate stores a new connection or updates the state of an exisitng
// Connection.
func ConnectionUpdate(c *controller.ConnectionController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadConnection{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		p.con.FromID = currentUser.ID

		if err := p.con.Validate(); err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		con, err := c.Update(app, p.con)
		if err != nil {
			if controller.IsInvalidEntity(err) {
				respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			} else {
				respondError(w, 0, err)
			}

			return
		}

		respondJSON(w, http.StatusOK, &payloadConnection{con: con})
	}
}

type payloadConnection struct {
	con *connection.Connection
}

func (p *payloadConnection) MarshalJSON() ([]byte, error) {
	f := struct {
		FromID       uint64    `json:"user_from_id"`
		FromIDString string    `json:"user_from_id_string"`
		ToID         uint64    `json:"user_to_id"`
		ToIDString   string    `json:"user_to_id_string"`
		State        string    `json:"state"`
		Type         string    `json:"type"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{
		FromID:    p.con.FromID,
		ToID:      p.con.ToID,
		State:     string(p.con.State),
		Type:      string(p.con.Type),
		CreatedAt: p.con.CreatedAt,
		UpdatedAt: p.con.UpdatedAt,
	}

	f.FromIDString = strconv.FormatUint(p.con.FromID, 10)
	f.ToIDString = strconv.FormatUint(p.con.ToID, 10)

	return json.Marshal(&f)
}

func (p *payloadConnection) UnmarshalJSON(raw []byte) error {
	f := struct {
		ToID       uint64 `json:"user_to_id"`
		ToIDString string `json:"user_to_id_string"`
		State      string `json:"state"`
		Type       string `json:"type"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.con = &connection.Connection{
		ToID:  f.ToID,
		State: connection.State(f.State),
		Type:  connection.Type(f.Type),
	}

	if f.ToID == 0 || f.ToIDString != "" {
		id, err := strconv.ParseUint(f.ToIDString, 10, 64)
		if err != nil {
			return err
		}

		p.con.ToID = id
	}

	return nil
}

type payloadConnections struct {
	cons    connection.List
	origin  uint64
	userMap user.Map
}

func (p *payloadConnections) MarshalJSON() ([]byte, error) {
	f := struct {
		Incoming      []*payloadConnection `json:"incoming"`
		IncomingCount int                  `json:"incoming_connections_count"`
		Outgoing      []*payloadConnection `json:"outgoing"`
		OutgoingCount int                  `json:"outgoing_connections_count"`
		Users         payloadUserMap       `json:"users"`
		UsersCount    int                  `json:"users_count"`
	}{
		Incoming:   []*payloadConnection{},
		Outgoing:   []*payloadConnection{},
		Users:      mapUserPresentation(p.userMap),
		UsersCount: len(p.userMap),
	}

	for _, c := range p.cons {
		if c.FromID == p.origin {
			f.Outgoing = append(f.Outgoing, &payloadConnection{con: c})
		} else {
			f.Incoming = append(f.Incoming, &payloadConnection{con: c})
		}
	}

	f.IncomingCount = len(f.Incoming)
	f.OutgoingCount = len(f.Outgoing)

	return json.Marshal(f)
}

type payloadSocial struct {
	ConnectionIDs []string
	Platform      string
	State         connection.State
	Type          connection.Type
}

func (p *payloadSocial) UnmarshalJSON(raw []byte) error {
	f := struct {
		ConnectionIDs []string `json:"connection_ids"`
		Platform      string   `json:"platform"`
		State         string   `json:"state"`
		Type          string   `json:"type"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	if f.State != "" {
		s := connection.State(f.State)
		switch s {
		case connection.StatePending, connection.StateConfirmed, connection.StateRejected:
			p.State = s
		default:
			return fmt.Errorf("invalid state %s", f.State)
		}
	} else {
		p.State = connection.StateConfirmed
	}

	t := connection.Type(f.Type)

	switch t {
	case connection.TypeFollow, connection.TypeFriend:
		p.Type = t
	default:
		return fmt.Errorf("invalid type %s", f.Type)
	}

	p.ConnectionIDs = f.ConnectionIDs
	p.Platform = f.Platform

	return nil
}
