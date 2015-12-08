package object

import (
	"time"

	"github.com/tapglue/multiverse/platform/flake"
)

type memService struct {
	objects map[string]map[uint64]*Object
}

// NewMemService returns a memory backed implementation of Service.
func NewMemService() Service {
	return &memService{
		objects: map[string]map[uint64]*Object{},
	}
}

func (s *memService) Put(namespace string, object *Object) (*Object, error) {
	if err := object.Validate(); err != nil {
		return nil, err
	}

	bucket, ok := s.objects[namespace]
	if !ok {
		return nil, ErrNamespaceNotFound
	}

	if object.ObjectID != 0 {
		keep := false
		for _, o := range bucket {
			if o.ID == object.ObjectID {
				keep = true
			}
		}

		if !keep {
			return nil, ErrMissingReference
		}
	}

	if object.ID == 0 {
		id, err := flake.NextID(flakeNamespace(namespace))
		if err != nil {
			return nil, err
		}

		object.CreatedAt = time.Now()
		object.ID = id
	} else {
		keep := false

		for _, o := range bucket {
			if o.ID == object.ID {
				keep = true
				object.CreatedAt = o.CreatedAt
			}
		}

		if !keep {
			return nil, ErrNotFound
		}
	}

	object.UpdatedAt = time.Now()
	bucket[object.ID] = copy(object)

	return copy(object), nil
}

func (s *memService) Query(namespace string, opts QueryOptions) ([]*Object, error) {
	bucket, ok := s.objects[namespace]
	if !ok {
		return nil, ErrNamespaceNotFound
	}

	os := []*Object{}

	for id, object := range bucket {
		if object.Deleted != opts.Deleted {
			continue
		}

		if opts.Owned != nil {
			if object.Owned != *opts.Owned {
				continue
			}
		}

		if opts.ID != nil && id != *opts.ID {
			continue
		}

		if !inIDs(object.OwnerID, opts.OwnerIDs) {
			continue
		}

		if !inIDs(object.ObjectID, opts.ObjectIDs) {
			continue
		}

		if !inTypes(object.Type, opts.Types) {
			continue
		}

		if !inVisibilities(object.Visibility, opts.Visibilities) {
			continue
		}

		os = append(os, object)
	}

	return os, nil
}

func (s *memService) Remove(namespace string, id uint64) error {
	bucket, ok := s.objects[namespace]
	if !ok {
		return ErrNamespaceNotFound
	}

	delete(bucket, id)

	return nil
}

func (s *memService) Setup(namespace string) error {
	if _, ok := s.objects[namespace]; !ok {
		s.objects[namespace] = map[uint64]*Object{}
	}

	return nil
}

func (s *memService) Teardown(namespace string) error {
	if _, ok := s.objects[namespace]; ok {
		delete(s.objects, namespace)
	}

	return nil
}

func copy(o *Object) *Object {
	old := *o
	return &old
}

func inIDs(id uint64, ids []uint64) bool {
	if len(ids) == 0 {
		return true
	}

	keep := false

	for _, i := range ids {
		if id == i {
			keep = true
			break
		}
	}

	return keep
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

func inVisibilities(visibility Visibility, vs []Visibility) bool {
	if len(vs) == 0 {
		return true
	}

	keep := false

	for _, v := range vs {
		if visibility == v {
			keep = true
			break
		}
	}

	return keep
}
