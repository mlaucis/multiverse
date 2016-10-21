package member

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/multiverse/platform/pg"

	"github.com/jmoiron/sqlx"
)

const (
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
	} else {
		query = pgInsertMember
	}

	input.UpdatedAt = now

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

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
	return nil, fmt.Errorf("Query not implemented")
}

func (s *pgService) Setup(ns string) error {
	return fmt.Errorf("Setup not implemented")
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

type pgSessionService struct {
	db *sqlx.DB
}

func PostgresSessionService(db *sqlx.DB) SessionService {
	return &pgSessionService{}
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
