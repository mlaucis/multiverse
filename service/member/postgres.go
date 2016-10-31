package member

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/multiverse/platform/pg"

	"github.com/jmoiron/sqlx"
)

const (
	pgInsertMember = `
		INSERT INTO
			%s.account_users(account_id, json_data)
		VALUES($1, $2)
		RETURNING id`
	pgUpdateMember = `
		UPDATE
			%s.account_users
		SET
			json_data = $3
		WHERE
			account_id = $1
			AND id = $2
		RETURNING id`

	pgListMembers = `
		SELECT
			id, account_id, json_data
		FROM
			%s.account_users
		%s`

	pgClauseEnabled = `(json_data->>'enabled')::BOOL = ?::BOOL`
	pgClauseIDs     = `id IN (?)`
	pgClauseOrgIDs  = `account_id IN (?)`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.account_users(
		id SERIAL PRIMARY KEY NOT NULL,
		account_id INT NOT NULL,
		json_data JSONB NOT NULL
	)`
	pgDropTable = `DROP TABLE IF EXISTS %s.account_users`
)

type pgService struct {
	db *sqlx.DB
}

// PostgresService returns a Postgres based Service implementation.
func PostgresService(db *sqlx.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Put(ns string, input *Member) (*Member, error) {
	var (
		now    = time.Now().UTC()
		query  = pgUpdateMember
		params = []interface{}{input.OrgID}
	)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if input.ID != 0 {
		params = append(params, input.ID)

		where, params, err := convertOpts(QueryOpts{
			IDs: []uint64{
				input.ID,
			},
		})
		if err != nil {
			return nil, err
		}

		ms, err := s.listMembers(ns, where, params...)
		if err != nil {
			return nil, err
		}

		if len(ms) != 1 {
			return nil, ErrNotFound
		}

		input.CreatedAt = ms[0].CreatedAt.UTC()
	} else {
		query = pgInsertMember

		if input.CreatedAt.IsZero() {
			input.CreatedAt = now
		}

		input.CreatedAt = input.CreatedAt.UTC()
	}

	input.LastLogin = input.LastLogin.UTC()
	input.UpdatedAt = now.UTC()

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	query = pg.WrapNamespace(query, ns)
	params = append(params, data)

	var id uint64

	err = s.db.QueryRow(query, params...).Scan(&id)
	if err != nil {
		if pg.IsRelationNotFound(pg.WrapError(err)) {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}

		err = s.db.QueryRow(query, params...).Scan(&id)
	}

	input.ID = id

	return input, err
}

func (s *pgService) Query(ns string, opts QueryOpts) (List, error) {
	where, params, err := convertOpts(opts)
	if err != nil {
		return nil, err
	}

	return s.listMembers(ns, where, params...)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		pg.WrapNamespace(pgCreateSchema, ns),
		pg.WrapNamespace(pgCreateTable, ns),
	}

	for _, q := range qs {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("query faield (%s): %s", q, err)
		}
	}

	return nil
}

func (s *pgService) Teardown(ns string) error {
	qs := []string{
		pg.WrapNamespace(pgDropTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) listMembers(
	ns, where string,
	params ...interface{},
) (List, error) {
	query := fmt.Sprintf(pgListMembers, ns, where)

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

	ms := List{}

	for rows.Next() {
		var (
			member = &Member{}

			id, orgID uint64
			raw       []byte
		)

		err := rows.Scan(&id, &orgID, &raw)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, member)
		if err != nil {
			return nil, err
		}

		member.ID = id
		member.OrgID = orgID

		ms = append(ms, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

type pgSessionService struct {
	db *sqlx.DB
}

func PostgresSessionService(db *sqlx.DB) SessionService {
	return &pgSessionService{
		db: db,
	}
}

func (s *pgSessionService) Put(ns string, input *Session) (*Session, error) {
	return nil, fmt.Errorf("Put not implemented")
}

func (s *pgSessionService) Query(ns string, opts SessionQueryOpts) (SessionList, error) {
	return nil, fmt.Errorf("Query not implemented")
}

func (s *pgSessionService) Setup(ns string) error {
	return fmt.Errorf("Setup not implemented")
}

func (s *pgSessionService) Teardown(ns string) error {
	return fmt.Errorf("Teardown not implemented")
}

func convertOpts(opts QueryOpts) (string, []interface{}, error) {
	var (
		clauses = []string{}
		params  = []interface{}{}
	)

	if opts.Enabled != nil {
		clause, _, err := sqlx.In(pgClauseEnabled, []interface{}{*opts.Enabled})
		if err != nil {
			return "", nil, err
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
			return "", nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.OrgIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.OrgIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseOrgIDs, ps)
		if err != nil {
			return "", nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	query := ""

	if len(clauses) > 0 {
		query = sqlx.Rebind(sqlx.DOLLAR, pg.ClausesToWhere(clauses...))
	}

	return query, params, nil
}
