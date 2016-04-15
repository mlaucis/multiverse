package connection

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/platform/pg"
)

const (
	pgDeleteConnection = `DELETE
		FROM %s.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
		AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT
		AND (json_data->>'type')::TEXT = $3::TEXT`
	pgInsertConnection = `INSERT INTO %s.connections(json_data) VALUES($1)`
	pgUpdateConnection = `UPDATE %s.connections
		SET json_data = $4
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
		AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT
		AND (json_data->>'type')::TEXT = $3::TEXT`

	pgCountConnections = `SELECT count(json_data) FROM %s.connections
		%s`
	pgListConnections = `SELECT json_data
		FROM %s.connections
		%s`
	pgClauseEnabled = `(json_data->>'enabled')::BOOL = ?::BOOL`
	pgClauseFromIDs = `(json_data->>'user_from_id')::BIGINT IN (?)`
	pgClauseStates  = `(json_data->>'state')::TEXT IN (?)`
	pgClauseToIDs   = `(json_data->>'user_to_id')::BIGINT IN (?)`
	pgClauseTypes   = `(json_data->>'type')::TEXT IN (?)`

	pgOrderCreatedAt = `ORDER BY (json_data->>'created_at')::TIMESTAMP DESC`

	pgCreatedByDay = `SELECT count(*), to_date(json_data->>'created_at', 'YYYY-MM-DD') as bucket
		FROM %s.connections
		WHERE (json_data->>'created_at')::DATE >= '%s'
		AND (json_data->>'created_at')::DATE <= '%s'
		GROUP BY bucket
		ORDER BY bucket`

	pgIndexFromID = `CREATE INDEX %s ON %s.connections
		USING btree (((json_data->>'user_from_id')::BIGINT))`
	pgIndexToID = `CREATE INDEX %s ON %s.connections
		USING btree (((json_data->>'user_to_id')::BIGINT))`
	pgIndexType = `CREATE INDEX %s ON %s.connections
		USING btree (((json_data->>'type')::TEXT))`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.connections
		(json_data JSONB NOT NULL)`
	pgDropTable = `DROP TABLE IF EXISTS %s.connections`
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

	return s.countConnections(ns, clauses, params...)
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

func (s *pgService) Put(ns string, con *Connection) (*Connection, error) {
	if err := con.Validate(); err != nil {
		return nil, err
	}

	var (
		params = []interface{}{con.FromID, con.ToID, string(con.Type)}

		query string
	)

	cs, err := s.listConnections(
		ns,
		[]string{pgClauseFromIDs, pgClauseToIDs, pgClauseTypes},
		params...,
	)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		query = wrapNamespace(pgUpdateConnection, ns)

		con.CreatedAt = cs[0].CreatedAt
	} else {
		params = []interface{}{}
		query = wrapNamespace(pgInsertConnection, ns)

		con.CreatedAt = time.Now().UTC()
	}

	con.UpdatedAt = time.Now().UTC()

	data, err := json.Marshal(con)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(query, append(params, data)...)
	if err != nil {
		return nil, err
	}

	return con, nil
}

func (s *pgService) Query(ns string, opts QueryOptions) (List, error) {
	clauses, params, err := convertOpts(opts)
	if err != nil {
		return nil, err
	}

	return s.listConnections(ns, clauses, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
		pg.GuardIndex(ns, "connection_from", pgIndexFromID),
		pg.GuardIndex(ns, "connection_to", pgIndexToID),
		pg.GuardIndex(ns, "connection_type", pgIndexType),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s: %s", query, err)
		}
	}

	return nil
}

func (s *pgService) Teardown(ns string) error {
	_, err := s.db.Exec(wrapNamespace(pgDropTable, ns))
	return err
}

func (s *pgService) countConnections(
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
		sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(pgCountConnections, ns, c)),
		params...,
	)
	if err != nil && pg.IsRelationNotFound(pg.WrapError(err)) {
		if err := s.Setup(ns); err != nil {
			return 0, err
		}

		err = s.db.Get(
			&count,
			sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(pgCountConnections, ns, c)),
			params...,
		)
	}

	return count, err
}

func (s *pgService) listConnections(
	ns string,
	clauses []string,
	params ...interface{},
) (List, error) {
	c := strings.Join(clauses, "\nAND")

	if len(clauses) > 0 {
		c = fmt.Sprintf("WHERE %s", c)
	}

	query := strings.Join([]string{
		fmt.Sprintf(pgListConnections, ns, c),
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
		}

		return nil, err
	}
	defer rows.Close()

	cs := List{}

	for rows.Next() {
		var (
			con = &Connection{}

			raw []byte
		)

		err := rows.Scan(&raw)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, con)
		if err != nil {
			return nil, err
		}

		cs = append(cs, con)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func convertOpts(opts QueryOptions) ([]string, []interface{}, error) {
	var (
		clauses = []string{}
		params  = []interface{}{}
	)

	if opts.Enabled != nil {
		clause, _, err := sqlx.In(pgClauseEnabled, []interface{}{*opts.Enabled})
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Enabled)
	}

	if len(opts.FromIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.FromIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseFromIDs, ps)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.States) > 0 {
		ps := []interface{}{}

		for _, state := range opts.States {
			ps = append(ps, string(state))
		}

		clause, _, err := sqlx.In(pgClauseStates, ps)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.ToIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.ToIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseToIDs, ps)
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.Types) > 0 {
		ps := []interface{}{}

		for _, t := range opts.Types {
			ps = append(ps, string(t))
		}

		clause, _, err := sqlx.In(pgClauseTypes, ps)
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
