package member

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type pgService struct {
	db *sqlx.DB
}

// PostgresService returns a Postgres based Service implementation.
func PostgresService(db *sqlx.DB) Service {
	return &pgService{}
}

func (s *pgService) Put(ns string, input *Member) (*Member, error) {
	return nil, fmt.Errorf("Put not implemented")
}

func (s *pgService) Query(ns string, opts QueryOpts) (List, error) {
	return nil, fmt.Errorf("Query not implemented")
}

func (s *pgService) Setup(ns string) error {
	return fmt.Errorf("Setup not implemented")
}

func (s *pgService) Teardown(ns string) error {
	return fmt.Errorf("Teardown not implemented")
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
