package org

import (
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/service"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// List is a collection of orgs.
type List []*Org

// Metadata is a bucket to provide custom org information.
type Metadata map[string]string

// Org is the root of the entire data model and holds apps and members.
type Org struct {
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	ID          uint64    `json:"-"`
	Metadata    Metadata  `json:"metadata"`
	Name        string    `json:"name"`
	PublicID    string    `json:"id"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QueryOptions is used to narrow-down user queries.
type QueryOptions struct {
	Enabled   *bool
	IDs       []uint64
	PublicIDs []string
	Tokens    []string
}

// Service for org interactions.
type Service interface {
	service.Lifecycle

	Put(*Org) (*Org, error)
	Query(QueryOptions) (List, error)
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type SsrviceMiddleware func(Service) Service

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FindByKey(string) (*v04_entity.Organization, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
