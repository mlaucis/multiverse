package connection

import (
	"fmt"
	"time"

	"github.com/tapglue/multiverse/platform/metrics"
)

type memService struct {
	cons map[string]map[string]*Connection
}

// NewMemService returns a memory backed implementation of Service.
func NewMemService() Service {
	return &memService{
		cons: map[string]map[string]*Connection{},
	}
}

func (s *memService) CreatedByDay(
	ns string,
	start, end time.Time,
) (metrics.Timeseries, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *memService) Put(ns string, con *Connection) (*Connection, error) {
	if err := s.Setup(ns); err != nil {
		return nil, err
	}

	if err := con.Validate(); err != nil {
		return nil, err
	}

	con.CreatedAt = time.Now().UTC()

	stored, ok := s.cons[ns][stringKey(con)]
	if ok {
		con.CreatedAt = stored.CreatedAt
	}

	con.UpdatedAt = time.Now().UTC()

	s.cons[ns][stringKey(con)] = con

	return con, nil
}

func (s *memService) Query(ns string, opts QueryOptions) (List, error) {
	if err := s.Setup(ns); err != nil {
		return nil, err
	}

	cs := List{}

	for _, con := range s.cons[ns] {
		if opts.Enabled != nil && con.Enabled != *opts.Enabled {
			continue
		}

		if !inIDs(con.FromID, opts.FromIDs) {
			continue
		}

		if !inStates(con.State, opts.States) {
			continue
		}

		if !inIDs(con.ToID, opts.ToIDs) {
			continue
		}

		if !inTypes(con.Type, opts.Types) {
			continue
		}

		cs = append(cs, con)
	}

	return cs, nil
}

func (s *memService) Setup(ns string) error {
	_, ok := s.cons[ns]
	if ok {
		return nil
	}

	s.cons[ns] = map[string]*Connection{}

	return nil
}

func (s *memService) Teardown(ns string) error {
	return fmt.Errorf("not implemented")
}

func inIDs(id uint64, ids []uint64) bool {
	if len(ids) == 0 {
		return true
	}

	keep := false

	for _, i := range ids {
		if i == id {
			keep = true
			break
		}
	}

	return keep
}

func inStates(s State, ss []State) bool {
	if len(ss) == 0 {
		return true
	}

	keep := false

	for _, state := range ss {
		if s == state {
			keep = true
			break
		}
	}

	return keep
}

func inTypes(t Type, ts []Type) bool {
	if len(ts) == 0 {
		return true
	}

	keep := false

	for _, ty := range ts {
		if t == ty {
			keep = true
			break
		}
	}

	return keep
}

func stringKey(con *Connection) string {
	return fmt.Sprintf("%d-%d-%s", con.FromID, con.ToID, con.Type)
}
