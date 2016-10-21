package member

import (
	"time"

	"github.com/tapglue/multiverse/platform/service"
)

// NamespaceDefault is the default namespace to isolate top-level data sets.
const NamespaceDefault = "tg"

// List is a Member collection.
type List []*Member

// Member is the representation of a user of an Org.
type Member struct {
	Email        string    `json:"email"`
	Enabled      bool      `json:"enabled"`
	Firstname    string    `json:"first_name"`
	LastLogin    time.Time `json:"last_login"`
	Lastname     string    `json:"last_name"`
	ID           uint64    `json:"-"`
	OrgID        uint64    `json:"-"`
	Password     string    `json:"password"`
	PublicID     string    `json:"id"`
	PublicOrgID  string    `json:"account_id"`
	SessionToken string    `json:"-"`
	URL          string    `json:"url"`
	Username     string    `json:"user_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// QueryOpts is used to narrow-down member queries.
type QueryOpts struct {
	Enabled *bool
	IDs     []uint64
}

// Service for member interactions.
type Service interface {
	service.Lifecycle

	Put(namespace string, member *Member) (*Member, error)
	Query(namespace string, opts QueryOpts) (List, error)
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// Session for Member authenticated interactions.
type Session struct {
	CreatedAt time.Time
	ID        string
	MemberID  uint64
	OrgID     uint64
}

// SessionList is a Session collection.
type SessionList []*Session

// SessionQueryOpts is used to narrow-down Session queries.
type SessionQueryOpts struct {
	IDs []string
}

// SessionService for Session interactions.
type SessionService interface {
	service.Lifecycle

	Put(namespace string, session *Session) (*Session, error)
	Query(namespace string, opts SessionQueryOpts) (SessionList, error)
}

// SessionServiceMiddleware is a chainbale behaviour modifier for SessionService.
type SessionServiceMiddleware func(SessionService) SessionService
