package event

import (
	"testing"
	"time"
)

func TestMemCount(t *testing.T) {
	testServiceCount(prepareMem, t)
}

func TestMemCreatedByDay(t *testing.T) {
	var (
		namespace = "created-by-day"
		service   = prepareMem(namespace, t)
		bucket    = map[uint64]*Event{}
	)

	for i := 0; i < 5; i++ {
		bucket[uint64(i)] = &Event{
			CreatedAt: time.Now().Add((time.Duration(i) * day) * -1),
			ID:        uint64(i),
		}
	}

	service.(*memService).events[namespace] = bucket

	ts, err := service.CreatedByDay(namespace, time.Now().Add(-3*day), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ts), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestMemPut(t *testing.T) {
	testServicePut(prepareMem, t)
}

func TestMemQuery(t *testing.T) {
	testServiceQuery(prepareMem, t)
}

func prepareMem(ns string, t *testing.T) Service {
	s := NewMemService()

	err := s.Teardown(ns)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Setup(ns)
	if err != nil {
		t.Fatal(err)
	}

	return s
}
