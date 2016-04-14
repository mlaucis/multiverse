package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/tapglue/multiverse/platform/pg"
)

const (
	pgClauseEnabled      = `(json_data->>'enabled')::BOOL = ?::BOOL`
	pgClauseInProduction = `(json_data->>'in_production')::BOOL = ?::BOOL`

	pgListApps = `SELECT id, account_id, json_data FROM tg.applications
		%s`

	pgOrderCreatedAt = `ORDER BY (json_data->>'created_at')::TIMESTAMP DESC`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.applications (
	  id SERIAL PRIMARY KEY NOT NULL,
	  account_id INT NOT NULL,
	  json_data JSONB NOT NULL,
	  enabled INT DEFAULT 1 NOT NULL
	)`
	pgDropTable = `DROP TABLE IF EXISTS %s.applications`

	pgCreateIndexCreatedAt = `CREATE INDEX %s ON %s.applications
		USING btree (((json_data->>'created_at')::TIMESTAMP))`
	pgCreateIndexBackendToken = `CREATE INDEX %s ON %s.applications
		USING BTREE (((json_data->>'backend_token')::TEXT))`
	pgCreateIndexToken = `CREATE INDEX %s ON %s.applications
		USING BTREE (((json_data->>'token')::TEXT))`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{db: db}
}

func (s *pgService) Query(opts QueryOptions) (List, error) {
	clauses, params, err := convertOpts(opts)
	if err != nil {
		return nil, err
	}

	return s.listApps(clauses, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
		pg.GuardIndex(ns, "apps_created_at", pgCreateIndexCreatedAt),
		pg.GuardIndex(ns, "apps_backend_token", pgCreateIndexBackendToken),
		pg.GuardIndex(ns, "apps_token", pgCreateIndexToken),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) listApps(clauses []string, params ...interface{}) (List, error) {
	c := strings.Join(clauses, "\nAND")

	if len(clauses) > 0 {
		c = fmt.Sprintf("WHERE %s", c)
	}

	query := strings.Join([]string{
		fmt.Sprintf(pgListApps, c),
		pgOrderCreatedAt,
	}, "\n")

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(pg.MetaNamespace); err != nil {
				return nil, err
			}

			rows, err = s.db.Query(query, params...)
			if err != nil {
				return nil, err
			}
		}
	}
	defer rows.Close()

	as := List{}

	for rows.Next() {
		var (
			app = &App{}

			id, orgID uint64
			raw       []byte
		)

		err := rows.Scan(&id, &orgID, &raw)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, app)
		if err != nil {
			return nil, err
		}

		app.ID = id
		app.OrgID = orgID

		as = append(as, app)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return as, nil
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

	if opts.InProduction != nil {
		clause, _, err := sqlx.In(pgClauseInProduction, []interface{}{*opts.InProduction})
		if err != nil {
			return nil, nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Enabled)
	}

	return clauses, params, nil
}

func wrapNamespace(query, namespace string) string {
	return fmt.Sprintf(query, namespace)
}
