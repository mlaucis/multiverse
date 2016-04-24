package user

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/platform/pg"
)

const (
	pgInserUser      = `INSERT INTO %s.users(json_data) VALUES($1)`
	pgUpdateLastRead = `
		UPDATE
			%s.users
		SET
			last_read = $2
		WHERE
			(json_data->>'id')::BIGINT = $1::BIGINT AND
			(json_data->>'enabled')::BOOL = true`
	pgUpdateUser = `
		UPDATE
			%s.users
		SET
			json_data = $1
		WHERE
			(json_data->>'id')::BIGINT = $2::BIGINT`

	pgClauseCustomIDs = `(json_data->>'custom_id')::TEXT IN (?)`
	pgClauseDeleted   = `(json_data->>'deleted')::BOOL = ?::BOOL`
	pgClauseEnabled   = `(json_data->>'enabled')::BOOL = ?::BOOL`
	pgClauseIDs       = `(json_data->>'id')::BIGINT IN (?)`

	pgOrderCreatedAt = `ORDER BY json_data->>'created_at' DESC`

	pgCountUsers = `SELECT count(json_data) FROM %s.users
		%s`
	pgListEvents = `SELECT json_data, last_read FROM %s.users
		%s`

	pgCreatedByDay = `SELECT count(*), to_date(json_data->>'created_at', 'YYYY-MM-DD') as bucket
		FROM %s.users
		WHERE (json_data->>'created_at')::DATE >= '%s'
		AND (json_data->>'created_at')::DATE <= '%s'
		GROUP BY bucket
		ORDER BY bucket`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.users (
		json_data JSONB NOT NULL,
		last_read TIMESTAMP DEFAULT '0001-01-01 00:00:00 UTC' NOT NULL
	)`
	pgDropTable = `DROP TABLE IF EXISTS %s.users`

	pgIndexCustomID = `CREATE INDEX %s ON %s.users
		USING btree (((json_data->>'custom_id')::TEXT))`
	pgIndexID = `CREATE INDEX %s ON %s.users
		USING btree (((json_data->>'id')::BIGINT))`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{db: db}
}

func (s *pgService) Count(ns string, opts QueryOptions) (int, error) {
	clauses, params, err := convertOpts(opts)
	if err != nil {
		return 0, err
	}

	return s.countUsers(ns, clauses, params...)
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
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}

			rows, err = s.db.Query(query)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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

func (s *pgService) Put(ns string, user *User) (*User, error) {
	var (
		now   = time.Now().UTC()
		query = pgUpdateUser

		params []interface{}
	)

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if user.ID != 0 {
		params = []interface{}{
			user.ID,
		}

		us, err := s.Query(ns, QueryOptions{
			IDs: []uint64{
				user.ID,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(us) == 0 {
			return nil, ErrNotFound
		}

		user.CreatedAt = us[0].CreatedAt
	} else {
		id, err := flake.NextID(flakeNamespace(ns))
		if err != nil {
			return nil, err
		}

		if user.CreatedAt.IsZero() {
			user.CreatedAt = now
		}
		user.ID = id
		user.LastRead = user.LastRead.UTC()

		query = pgInserUser
	}

	user.UpdatedAt = now

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	params = append([]interface{}{data}, params...)

	_, err = s.db.Exec(wrapNamespace(query, ns), params...)
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}

			_, err = s.db.Exec(wrapNamespace(query, ns), params...)
		}
	}

	return user, err
}

func (s *pgService) PutLastRead(ns string, userID uint64, ts time.Time) error {
	_, err := s.db.Exec(wrapNamespace(pgUpdateLastRead, ns), userID, ts.UTC())
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return err
			}

			_, err = s.db.Exec(pgUpdateLastRead, userID, ts.UTC())
		}
	}
	return err
}

func (s *pgService) Query(ns string, opts QueryOptions) (List, error) {
	clauses, params, err := convertOpts(opts)
	if err != nil {
		return nil, err
	}

	return s.listUsers(ns, clauses, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
		pg.GuardIndex(ns, "user_custom_id", pgIndexCustomID),
		pg.GuardIndex(ns, "user_id", pgIndexID),
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

func (s *pgService) countUsers(
	ns string,
	clauses []string,
	params ...interface{},
) (int, error) {
	c := strings.Join(clauses, "\nAND")

	if len(clauses) > 0 {
		c = fmt.Sprintf("WHERE %s", c)
	}

	count := 0

	err := s.db.Get(
		&count,
		sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(pgCountUsers, ns, c)),
		params...,
	)
	if err != nil && pg.IsRelationNotFound(pg.WrapError(err)) {
		if err := s.Setup(ns); err != nil {
			return 0, err
		}

		err = s.db.Get(
			&count,
			sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(pgCountUsers, ns, c)),
			params...,
		)
	}

	return count, err
}

func (s *pgService) listUsers(
	ns string,
	clauses []string,
	params ...interface{},
) (List, error) {
	c := strings.Join(clauses, "\nAND")

	if len(clauses) > 0 {
		c = fmt.Sprintf("WHERE %s", c)
	}

	query := strings.Join([]string{
		fmt.Sprintf(pgListEvents, ns, c),
		pgOrderCreatedAt,
	}, "\n")

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}

			rows, err = s.db.Query(query, params...)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	us := List{}

	for rows.Next() {
		var (
			user = &User{}

			lastRead time.Time
			raw      []byte
		)

		if err := rows.Scan(&raw, &lastRead); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(raw, user); err != nil {
			return nil, err
		}

		user.LastRead = lastRead.UTC()

		us = append(us, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return us, nil
}

func convertOpts(opts QueryOptions) ([]string, []interface{}, error) {
	var (
		clauses = []string{}
		params  = []interface{}{}
	)

	if len(opts.CustomIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.CustomIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseCustomIDs, ps)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if opts.Deleted != nil {
		clause, _, err := sqlx.In(pgClauseDeleted, []interface{}{*opts.Deleted})
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Deleted)
	}

	if opts.Enabled != nil {
		clause, _, err := sqlx.In(pgClauseEnabled, []interface{}{*opts.Enabled})
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Enabled)
	}

	if len(opts.IDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.IDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseIDs, ps)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	return clauses, params, nil
}

func wrapNamespace(query, namespace string) string {
	return fmt.Sprintf(query, namespace)
}
