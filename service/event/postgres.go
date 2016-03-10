package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/platform/metrics"
)

const (
	pgInsertEvent = `INSERT INTO %s.events(json_data) VALUES($1)`
	pgUpdateEvent = `UPDATE %s.events SET json_data = $1
		WHERE (json_data->>'id')::BIGINT = $2::BIGINT`
	pgDeleteEvent = `DELETE FROM %s.events
		WHERE (json_data->>'id')::BIGINT = $1::BIGINT`

	pgSelectEvents = `SELECT json_data FROM %s.events
		%s`

	pgClauseEnabled             = `(json_data->>'enabled')::BOOL = ?::BOOL`
	pgClauseExternalObjectIDs   = `(json_data->'object'->>'id')::TEXT IN (?)`
	pgClauseExternalObjectTypes = `(json_data->'object'->>'type')::TEXT in (?)`
	pgClauseIDs                 = `(json_data->>'id')::BIGINT IN (?)`
	pgClauseObjectIDs           = `(json_data->>'object_id')::BIGINT IN (?)`
	pgClauseOwned               = `(json_data->>'owned')::BOOL = ?::BOOL`
	pgClauseTypes               = `(json_data->>'type')::TEXT in (?)`
	pgClauseUserIDs             = `(json_data->>'user_id')::BIGINT IN (?)`
	pgClauseVisibilities        = `(json_data->>'visibility')::INT IN (?)`

	pgActiveByPeriod = `SELECT (json_data ->> 'user_id')::BIGINT AS userid, COUNT(*)
    FROM %s.events
      WHERE
        %s
    GROUP BY userid
    ORDER BY COUNT DESC`
	pgClauseByDay   = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 day'`
	pgClauseByWeek  = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 week'`
	pgClauseByMonth = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 month'`

	pgOrderCreatedAt = `ORDER BY json_data->>'created_at' DESC`

	pgCreatedByDay = `SELECT count(*), to_date(json_data->>'created_at', 'YYYY-MM-DD') as bucket
		FROM %s.events
		WHERE (json_data->>'created_at')::DATE >= '%s'
		AND (json_data->>'created_at')::DATE <= '%s'
		GROUP BY bucket
		ORDER BY bucket`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.events
		(json_data JSONB NOT NULL)`
	pgDropTable = `DROP TABLE IF EXISTS %s.events`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{db: db}
}

func (s *pgService) ActiveUserIDs(
	ns string,
	p Period,
) ([]uint64, error) {
	var clause string

	switch p {
	case ByDay:
		clause = pgClauseByDay
	case ByWeek:
		clause = pgClauseByWeek
	case ByMonth:
		clause = pgClauseByMonth
	default:
		return nil, fmt.Errorf("period %s not supported", p)
	}

	query := fmt.Sprintf(pgActiveByPeriod, ns, clause)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []uint64{}
	for rows.Next() {
		var (
			id    uint64
			count int
		)

		err := rows.Scan(&id, &count)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *pgService) CreatedByDay(
	ns string,
	start, end time.Time,
) (metrics.Timeseries, error) {
	query := fmt.Sprintf(
		pgCreatedByDay,
		ns,
		start.Format(metrics.BucketFormat),
		end.Format(metrics.BucketFormat),
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := []metrics.Datapoint{}
	for rows.Next() {
		var (
			bucket time.Time
			value  int
		)

		err := rows.Scan(&value, &bucket)
		if err != nil {
			return nil, err
		}

		ts = append(
			ts,
			metrics.Datapoint{
				Bucket: bucket.Format(metrics.BucketFormat),
				Value:  value,
			},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *pgService) Put(ns string, event *Event) (*Event, error) {
	var (
		query = pgUpdateEvent

		params []interface{}
	)

	if err := event.Validate(); err != nil {
		return nil, err
	}

	if event.ID != 0 {
		params = []interface{}{
			event.ID,
		}

		es, err := s.Query(ns, QueryOptions{
			IDs: []uint64{
				event.ID,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(es) == 0 {
			return nil, ErrNotFound
		}

		event.CreatedAt = es[0].CreatedAt
	} else {
		id, err := flake.NextID(flakeNamespace(ns))
		if err != nil {
			return nil, err
		}

		event.CreatedAt = time.Now().UTC()
		event.ID = id

		query = pgInsertEvent
	}

	event.UpdatedAt = time.Now().UTC()

	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	params = append([]interface{}{data}, params...)

	_, err = s.db.Exec(wrapNamespace(query, ns), params...)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *pgService) Query(ns string, opts QueryOptions) (List, error) {
	var (
		clauses = []string{}
		params  = []interface{}{}
	)

	if opts.Enabled != nil {
		clause, _, err := sqlx.In(pgClauseEnabled, []interface{}{*opts.Enabled})
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Enabled)
	}

	if len(opts.ExternalObjectIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.ExternalObjectIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseExternalObjectIDs, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.ExternalObjectTypes) > 0 {
		ps := []interface{}{}

		for _, t := range opts.ExternalObjectTypes {
			ps = append(ps, t)
		}

		clause, _, err := sqlx.In(pgClauseExternalObjectTypes, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.IDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.IDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseIDs, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.ObjectIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.ObjectIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseObjectIDs, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if opts.Owned != nil {
		clause, _, err := sqlx.In(pgClauseOwned, []interface{}{*opts.Owned})
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Owned)
	}

	if len(opts.Types) > 0 {
		ps := []interface{}{}

		for _, t := range opts.Types {
			ps = append(ps, t)
		}

		clause, _, err := sqlx.In(pgClauseTypes, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.UserIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.UserIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseUserIDs, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.Visibilities) > 0 {
		ps := []interface{}{}

		for _, v := range opts.Visibilities {
			ps = append(ps, v)
		}

		clause, _, err := sqlx.In(pgClauseVisibilities, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	c := strings.Join(clauses, "\nAND ")

	if len(clauses) > 0 {
		c = strings.Join([]string{
			"WHERE",
			c,
		}, " ")
	}

	query := fmt.Sprintf(
		pgSelectEvents,
		ns,
		c,
	)

	query = strings.Join([]string{
		query,
		pgOrderCreatedAt,
	}, "\n")

	return s.queryEvents(query, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) Teardown(ns string) error {
	qs := []string{
		wrapNamespace(pgDropTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) queryEvents(query string, params ...interface{}) (List, error) {
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, pgWrapError(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	es := List{}

	for rows.Next() {
		var (
			event = &Event{}

			raw []byte
		)

		err := rows.Scan(&raw)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, event)
		if err != nil {
			return nil, err
		}

		es = append(es, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return es, nil
}

func pgWrapError(err error) error {
	if err, ok := err.(*pq.Error); ok && err.Code == "42P01" {
		return ErrNamespaceNotFound
	}

	return err
}

func wrapNamespace(query, namespace string) string {
	return fmt.Sprintf(query, namespace)
}
