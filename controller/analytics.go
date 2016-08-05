package controller

import (
	"time"

	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/service/app"
)

// AnalyticsWhere combines query clauses for analytcs requests.
type AnalyticsWhere struct {
	End   time.Time
	Start time.Time
}

// AppResult bundles the entity timeseries.
type AppResult map[string]metrics.Timeseries

// Summary is the sums for the enity timeseries.
type Summary struct {
	NewConnections int
	NewEvents      int
	NewObjects     int
	NewUsers       int
}

// AnalyticsController bundles the business constraints for analytics endpoints
// for organisations.
type AnalyticsController struct {
	apps        app.Service
	connections metrics.BucketByDay
	events      metrics.BucketByDay
	objects     metrics.BucketByDay
	users       metrics.BucketByDay
}

// NewAnalyticsController returns a controller instance.
func NewAnalyticsController(
	apps app.Service,
	connections, events, objects, users metrics.BucketByDay,
) *AnalyticsController {
	return &AnalyticsController{
		apps:        apps,
		connections: connections,
		events:      events,
		objects:     objects,
		users:       users,
	}
}

// App returns the timeseries data for all entities of an app.
func (c *AnalyticsController) App(
	publicID string,
	where *AnalyticsWhere,
) (AppResult, error) {
	as, err := c.apps.Query(app.NamespaceDefault, app.QueryOptions{
		Enabled: &defaultEnabled,
		PublicIDs: []string{
			publicID,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(as) != 1 {
		return nil, ErrNotFound
	}

	app := as[0]

	if where == nil {
		where = &AnalyticsWhere{
			Start: time.Now().AddDate(0, -1, 0),
			End:   time.Now(),
		}
	}

	start, end := where.Start, where.End

	cs, err := c.connections.CreatedByDay(app.Namespace(), start, end)
	if err != nil {
		return nil, err
	}

	es, err := c.events.CreatedByDay(app.Namespace(), start, end)
	if err != nil {
		return nil, err
	}

	os, err := c.objects.CreatedByDay(app.Namespace(), start, end)
	if err != nil {
		return nil, err
	}

	us, err := c.users.CreatedByDay(app.Namespace(), start, end)
	if err != nil {
		return nil, err
	}

	return AppResult{
		"connections": cs,
		"events":      es,
		"objects":     os,
		"users":       us,
	}, nil
}
