package app

import (
	"fmt"

	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// App represents an Org owned data container.
type App struct {
	ID           uint64 `json:"-"`
	InProduction bool   `json:"in_production"`
	Name         string `json:"name"`
	OrgID        uint64 `json:"-"`
	// Missing fields
}

func (a *App) Namespace() string {
	return fmt.Sprintf("app_%d_%d", a.OrgID, a.ID)
}

// List is an App collection.
type List []*App

// QueryOptions are used to narrow down app queries.
type QueryOptions struct {
	Enabled      *bool
	InProduction *bool
}

// Service for app interactions.
type Service interface {
	Query(QueryOptions) (List, error)
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FindByApplicationToken(token string) (*v04_entity.Application, []errors.Error)
	FindByBackendToken(token string) (*v04_entity.Application, []errors.Error)
	FindByPublicID(publicID string) (*v04_entity.Application, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
