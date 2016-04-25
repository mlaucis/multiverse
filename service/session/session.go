package session

import (
	"time"

	"github.com/tapglue/multiverse/platform/service"
)

// List is a collection of sessions.
type List []*Session

// Map is a session collection with their id as index.
type Map map[string]*Session

// QueryOptions is used to narrow-down session queries.
type QueryOptions struct {
	Enabled *bool
	IDs     []string
}

// Service for session interactions
type Service interface {
	service.Lifecycle

	Put(namespace string, session *Session) (*Session, error)
	Query(namespace string, opts QueryOptions) (List, error)
}

// Session attaches a session id to a user id.
type Session struct {
	CreatedAt time.Time
	Enabled   bool
	ID        string
	UserID    uint64
}

// Validate performs semantic checks on the Session.
func (s *Session) Validate() error {
	if s.ID == "" {
		return wrapError(ErrInvalidSession, "id must be set")
	}

	if s.UserID == 0 {
		return wrapError(ErrInvalidSession, "UserID must be set")
	}

	return nil
}
