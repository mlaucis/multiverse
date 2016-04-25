package session

import "testing"

func TestMemPut(t *testing.T) {
	testServiecPut(t, prepareMem)
}

func TestMemQuery(t *testing.T) {
	testServiecQuery(t, prepareMem)
}

func prepareMem(t *testing.T, ns string) Service {
	s := NewMemService()

	if err := s.Teardown(ns); err != nil {
		t.Fatal(err)
	}

	return s
}
