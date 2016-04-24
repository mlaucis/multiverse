package user

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/platform/service"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	v04_errmsg "github.com/tapglue/multiverse/v04/errmsg"
)

// TargetType is the identifier used for events targeting a User.
const TargetType = "tg_user"

// Image represents a user image asset.
type Image struct {
	URL    string `json:"url"`
	Type   string `json:"type"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// List is a collection of users.
type List []*User

// Map is a user collection with their id as index.
type Map map[uint64]*User

// Metadata is a bucket to provide additional user information.
type Metadata map[string]string

// QueryOptions is used to narrow-down user queries.
type QueryOptions struct {
	CustomIDs []string
	Deleted   *bool
	Enabled   *bool
	IDs       []uint64
}

// Service for user interactions.
type Service interface {
	metrics.BucketByDay
	service.Lifecycle

	Count(namespace string, opts QueryOptions) (int, error)
	Put(namespace string, user *User) (*User, error)
	PutLastRead(namespace string, userID uint64, lastRead time.Time) error
	Query(namespace string, opts QueryOptions) (List, error)
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// User is the representation of a customer of an app.
type User struct {
	CustomID  string            `json:"custom_id,omitempty"`
	Deleted   bool              `json:"deleted"`
	Enabled   bool              `json:"enabled"`
	Email     string            `json:"email"`
	Firstname string            `json:"first_name"`
	ID        uint64            `json:"id"`
	Images    map[string]Image  `json:"images,omitempty"`
	Lastname  string            `json:"last_name"`
	LastRead  time.Time         `json:"-"`
	Metadata  Metadata          `json:"metadata"`
	Password  string            `json:"password"`
	SocialIDs map[string]string `json:"social_ids,omitempty"`
	URL       string            `json:"url,omitempty"`
	Username  string            `json:"user_name"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Validate performs semantic checks on the passed User values for correctness.
func (u *User) Validate() error {
	if u.Email == "" && u.Username == "" {
		return wrapError(ErrInvalidUser, "email or username must be set")
	}

	if ok := govalidator.IsEmail(u.Email); u.Email != "" && !ok {
		return wrapError(ErrInvalidUser, "invalid email address")
	}

	if u.Firstname != "" {
		if len(u.Firstname) < 2 {
			return wrapError(ErrInvalidUser, "firstname too short")
		}
		if len(u.Firstname) > 40 {
			return wrapError(ErrInvalidUser, "firstname too long")
		}
	}

	if u.Lastname != "" {
		if len(u.Lastname) < 2 {
			return wrapError(ErrInvalidUser, "lastname too short")
		}
		if len(u.Lastname) > 40 {
			return wrapError(ErrInvalidUser, "lastname too long")
		}
	}

	if ok := govalidator.IsURL(u.URL); u.URL != "" && !ok {
		return wrapError(ErrInvalidUser, "invalid url")
	}

	if u.Password == "" {
		return wrapError(ErrInvalidUser, "password must be set")
	}

	if u.Username != "" {
		if len(u.Username) < 2 {
			return wrapError(ErrInvalidUser, "username too short")
		}
		if len(u.Username) > 40 {
			return wrapError(ErrInvalidUser, "username too long")
		}
	}

	return nil
}

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FilterByEmail(
		orgID, appID int64,
		emails []string,
	) ([]*v04_entity.ApplicationUser, []errors.Error)
	FilterBySocialIDs(
		orgID, appID int64,
		platform string,
		ids []string,
	) ([]*v04_entity.ApplicationUser, []errors.Error)
	FindBySession(
		orgID, appID int64,
		key string,
	) (*v04_entity.ApplicationUser, []errors.Error)
	Read(
		orgID, appID int64,
		id uint64,
		stats bool,
	) (*v04_entity.ApplicationUser, []errors.Error)
	UpdateLastRead(orgID, appID int64, userID uint64) []errors.Error
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService

// StrangleMap is an old user collection with their id as index.
type StrangleMap map[uint64]*v04_entity.ApplicationUser

// Merge combines two StrangleMaps.
func (m StrangleMap) Merge(x StrangleMap) StrangleMap {
	for id, user := range x {
		m[id] = user
	}

	return m
}

// StrangleList is a collection of old users.
type StrangleList []*v04_entity.ApplicationUser

// IDs returns the list of user ids.
func (l StrangleList) IDs() []uint64 {
	ids := []uint64{}

	for _, user := range l {
		ids = append(ids, user.ID)
	}

	return ids
}

// ToStrangleMap turns the user list into a Map.
func (l StrangleList) ToStrangleMap() StrangleMap {
	m := StrangleMap{}

	for _, user := range l {
		m[user.ID] = user
	}

	return m
}

// StrangleMapFromIDs return a populated user map for the given list of ids.
func StrangleMapFromIDs(
	s StrangleService,
	app *v04_entity.Application,
	ids ...uint64,
) (StrangleMap, error) {
	um := StrangleMap{}

	for _, id := range ids {
		if _, ok := um[id]; ok {
			continue
		}

		user, errs := s.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == v04_errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		um[user.ID] = user
	}

	return um, nil
}

// StrangleListFromIDs gathers a user collection from the service for the given ids.
func StrangleListFromIDs(
	s StrangleService,
	app *v04_entity.Application,
	ids ...uint64,
) (StrangleList, error) {
	var (
		seen = map[uint64]struct{}{}
		us   = StrangleList{}
	)

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		u, errs := s.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == v04_errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		us = append(us, u)
	}

	return us, nil
}

func flakeNamespace(ns string) string {
	return fmt.Sprintf("%s_%s", ns, "users")
}
