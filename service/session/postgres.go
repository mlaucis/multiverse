package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/tapglue/multiverse/platform/pg"
)

const (
	pgInsertSession = `INSERT INTO
		%s.sessions(user_id, session_id, created_at, enabled)
		VALUES($1, $2, $3, $4)`

	pgClauseEnabled = `enabled = ?`
	pgClauseIDs     = `session_id IN (?)`

	pgOrderCreatedAt = `ORDER BY created_at DESC`

	pgListSessions = `
		SELECT
			user_id, session_id, created_at, enabled
		FROM
			%s.sessions
		%s`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.sessions (
		user_id BIGINT NOT NULL,
		session_id VARCHAR(40) NOT NULL,
		created_at TIMESTAMP DEFAULT now() NOT NULL,
		enabled BOOL DEFAULT TRUE NOT NULL
	)`
	pgDropTable = `DROP TABLE IF EXISTS %s.sessions`

	pgIndexID = `CREATE INDEX %s ON %s.sessions (session_id)`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{db: db}
}

func (s *pgService) Put(ns string, session *Session) (*Session, error) {
	if err := session.Validate(); err != nil {
		return nil, err
	}

	if session.CreatedAt.IsZero() {
		// FIXME: Postgres doesn't preserve nanosecond precision.
		format := "2006-01-02 15:04:05.000000 UTC"
		ts, err := time.Parse(format, time.Now().Format(format))
		if err != nil {
			return nil, err
		}
		session.CreatedAt = ts
	}

	session.CreatedAt = session.CreatedAt.UTC()

	var (
		query  = fmt.Sprintf(pgInsertSession, ns)
		params = []interface{}{
			session.UserID,
			session.ID,
			session.CreatedAt,
			session.Enabled,
		}
	)

	_, err := s.db.Exec(query, params...)
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}

			_, err = s.db.Exec(query, params...)
		}
	}

	return session, err
}

func (s *pgService) Query(ns string, opts QueryOptions) (List, error) {
	clauses, params, err := convertOpts(opts)
	if err != nil {
		return nil, err
	}

	return s.listSessions(ns, clauses, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		fmt.Sprintf(pgCreateSchema, ns),
		fmt.Sprintf(pgCreateTable, ns),
		pg.GuardIndex(ns, "session_id", pgIndexID),
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
		fmt.Sprintf(pgDropTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) listSessions(
	ns string,
	clauses []string,
	params ...interface{},
) (List, error) {
	c := strings.Join(clauses, "\nAND ")

	if len(clauses) > 0 {
		c = fmt.Sprintf("WHERE %s", c)
	}

	query := strings.Join([]string{
		fmt.Sprintf(pgListSessions, ns, c),
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

	ss := List{}

	for rows.Next() {
		s := &Session{}

		err := rows.Scan(
			&s.UserID,
			&s.ID,
			&s.CreatedAt,
			&s.Enabled,
		)
		if err != nil {
			return nil, err
		}

		s.CreatedAt = s.CreatedAt.UTC()

		ss = append(ss, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ss, nil
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
