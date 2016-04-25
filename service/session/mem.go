package session

import (
	"time"
)

type memService struct {
	sessions map[string]Map
}

// NewMemService returns a memory based Service implementation.
func NewMemService() Service {
	return &memService{
		sessions: map[string]Map{},
	}
}

func (s *memService) Put(ns string, session *Session) (*Session, error) {
	if err := s.Setup(ns); err != nil {
		return nil, err
	}

	if err := session.Validate(); err != nil {
		return nil, err
	}

	session.CreatedAt = time.Now().UTC()

	s.sessions[ns][session.ID] = copy(session)

	return copy(session), nil
}

func (s *memService) Query(ns string, opts QueryOptions) (List, error) {
	if err := s.Setup(ns); err != nil {
		return nil, err
	}

	return filterMap(s.sessions[ns], opts), nil
}

func (s *memService) Setup(ns string) error {
	if _, ok := s.sessions[ns]; !ok {
		s.sessions[ns] = Map{}
	}

	return nil
}

func (s *memService) Teardown(ns string) error {
	if _, ok := s.sessions[ns]; ok {
		delete(s.sessions, ns)
	}

	return nil
}

func copy(s *Session) *Session {
	old := *s
	return &old
}

func filterMap(sm Map, opts QueryOptions) List {
	ss := List{}

	for id, s := range sm {
		if opts.Enabled != nil && s.Enabled != *opts.Enabled {
			continue
		}

		if !inTypes(id, opts.IDs) {
			continue
		}

		ss = append(ss, s)
	}

	return ss
}

func inTypes(ty string, ts []string) bool {
	if len(ts) == 0 {
		return true
	}

	keep := false

	for _, t := range ts {
		if ty == t {
			keep = true
			break
		}
	}

	return keep
}
