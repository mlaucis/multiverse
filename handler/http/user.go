package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/user"
)

const userLocationFmt = "https://%s/%s/users/%d"

// UserCreate stores the provided user and returns it with a valid session.
func UserCreate(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp = appFromContext(ctx)
			p          = payloadUser{}
			version    = versionFromContext(ctx)
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		u, err := c.Create(currentApp, p.user)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		w.Header().Set("Location", fmt.Sprintf(
			userLocationFmt,
			r.Host,
			version,
			u.ID,
		))

		respondJSON(w, http.StatusCreated, &payloadUser{user: u})
	}
}

// UserDelete disbales the current user.
func UserDelete(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		err := c.Delete(currentApp, currentUser)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// UserLogin finds the user by email or username and creates a Session.
func UserLogin(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app = appFromContext(ctx)
			p   = payloadLogin{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		if p.email == "" && p.username == "" {
			respondError(w, 0, wrapError(ErrBadRequest, "email or user_name must be set"))
			return
		}

		u, err := c.LoginEmail(app, p.email, p.password)
		if err != nil {
			if !controller.IsNotFound(err) {
				respondError(w, 0, err)
				return
			}

			u, err = c.LoginUsername(app, p.username, p.password)
			if err != nil {
				respondError(w, 0, err)
				return
			}
		}

		respondJSON(w, http.StatusCreated, &payloadUser{user: u})
	}
}

// UserLogout finds the session of the user and destroys it.
func UserLogout(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			token       = tokenFromContext(ctx)
			tokenType   = tokenTypeFromContext(ctx)
		)

		if tokenType == tokenBackend {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		err := c.Logout(currentApp, currentUser.ID, token)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// UserRetrieve returns the user for the requested id.
func UserRetrieve(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		userID, err := strconv.ParseUint(mux.Vars(r)["userID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		u, err := c.Retrieve(currentApp, currentUser.ID, userID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUser{user: u})
	}
}

// UserRetrieveMe returns the current user.
func UserRetrieveMe(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		u, err := c.Retrieve(currentApp, currentUser.ID, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUser{user: u})
	}
}

// UserSearchEmails returns all Users for the emails of the payload.
func UserSearchEmails(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadSearchEmails{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		if len(p.Emails) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := c.ListByEmails(app, currentUser.ID, p.Emails...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{
			users: us,
		})
	}
}

// UserSearch returns all users for the given search query.
func UserSearch(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			query       = r.URL.Query().Get("q")
		)

		if len(query) < 3 {
			respondError(w, 0, wrapError(ErrBadRequest, "query must be over 3 characters"))
			return
		}

		us, err := c.Search(currentApp, currentUser.ID, query)
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

// UserSearchPlatform returns all users for the given ids and platform.
func UserSearchPlatform(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			platform    = mux.Vars(r)["platform"]
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadSearchPlatform{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		if len(p.IDs) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := c.ListByPlatformIDs(app, currentUser.ID, platform, p.IDs...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{
			users: us,
		})
	}
}

// UserUpdate stores the new attributes given.
func UserUpdate(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadUser{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		u, err := c.Update(currentApp, currentUser, p.user)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUser{user: u})
	}
}

type payloadLogin struct {
	email    string
	password string
	username string
	wildcard string
}

func (p *payloadLogin) UnmarshalJSON(raw []byte) error {
	f := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"user_name"`
		Wildcard string `json:"username"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	if f.Password == "" {
		return fmt.Errorf("password must be set")
	}

	if f.Wildcard != "" {
		f.Email, f.Username = f.Wildcard, f.Wildcard
	}

	if f.Email == "" && f.Username == "" {
		return fmt.Errorf("email or user_name must be provided")
	}

	p.email = f.Email
	p.password = f.Password
	p.username = f.Username

	return nil
}

type payloadSearchEmails struct {
	Emails []string `json:"emails"`
}

type payloadSearchPlatform struct {
	IDs []string `json:"ids"`
}

type payloadUser struct {
	user *user.User
}

func (p *payloadUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		CustomID       string                `json:"custom_id,omitempty"`
		Email          string                `json:"email"`
		Firstname      string                `json:"first_name"`
		FollowerCount  int                   `json:"follower_count"`
		FollowingCount int                   `json:"followed_count"`
		FriendCount    int                   `json:"friend_count"`
		ID             uint64                `json:"id"`
		IDString       string                `json:"id_string"`
		Images         map[string]user.Image `json:"images,omitempty"`
		IsFollower     bool                  `json:"is_follower"`
		IsFollowing    bool                  `json:"is_followed"`
		IsFriend       bool                  `json:"is_friend"`
		Lastname       string                `json:"last_name"`
		Metadata       user.Metadata         `json:"metadata"`
		SessionToken   string                `json:"session_token,omitempty"`
		SocialIDs      map[string]string     `json:"social_ids,omitempty"`
		URL            string                `json:"url,omitempty"`
		Username       string                `json:"user_name"`
		CreatedAt      time.Time             `json:"created_at"`
		UpdatedAt      time.Time             `json:"updated_at"`
	}{
		CustomID:       p.user.CustomID,
		Email:          p.user.Email,
		Firstname:      p.user.Firstname,
		FollowerCount:  p.user.FollowerCount,
		FollowingCount: p.user.FollowingCount,
		FriendCount:    p.user.FriendCount,
		ID:             p.user.ID,
		IDString:       strconv.FormatUint(p.user.ID, 10),
		Images:         p.user.Images,
		IsFollower:     p.user.IsFollower,
		IsFollowing:    p.user.IsFollowing,
		IsFriend:       p.user.IsFriend,
		Lastname:       p.user.Lastname,
		Metadata:       p.user.Metadata,
		SessionToken:   p.user.SessionToken,
		SocialIDs:      p.user.SocialIDs,
		URL:            p.user.URL,
		Username:       p.user.Username,
		CreatedAt:      p.user.CreatedAt,
		UpdatedAt:      p.user.UpdatedAt,
	})
}

func (p *payloadUser) UnmarshalJSON(raw []byte) error {
	f := struct {
		CustomID  string                `json:"custom_id,omitempty"`
		Email     string                `json:"email"`
		Firstname string                `json:"first_name"`
		Images    map[string]user.Image `json:"images,omitempty"`
		Lastname  string                `json:"last_name"`
		Metadata  user.Metadata         `json:"metadata"`
		Password  string                `json:"password,omitempty"`
		SocialIDs map[string]string     `json:"social_ids"`
		URL       string                `json:"url,omitempty"`
		Username  string                `json:"user_name"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.user = &user.User{
		CustomID:  f.CustomID,
		Email:     f.Email,
		Firstname: f.Firstname,
		Images:    f.Images,
		Lastname:  f.Lastname,
		Metadata:  f.Metadata,
		Password:  f.Password,
		SocialIDs: f.SocialIDs,
		URL:       f.URL,
		Username:  f.Username,
	}

	return p.user.Validate()
}

type payloadUserMap struct {
	userMap user.Map
}

func (p *payloadUserMap) MarshalJSON() ([]byte, error) {
	m := map[string]*payloadUser{}

	for id, u := range p.userMap {
		m[strconv.FormatUint(id, 10)] = &payloadUser{user: u}
	}

	return json.Marshal(m)
}
