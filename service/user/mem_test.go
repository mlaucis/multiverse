package user

import "testing"

func TestMemCount(t *testing.T) {
	testServiceCount(t, prepareMem)
}

func TestMemCreatedByDay(t *testing.T) {
	testServiceCreatedByDay(t, prepareMem)
}

func TestMemPut(t *testing.T) {
	testServicePut(t, prepareMem)
}

func TestMemPutLastRead(t *testing.T) {
	testServicePutLastRead(t, prepareMem)
}

func TestMemSearch(t *testing.T) {
	testServiceSearch(t, prepareMem)
}

func prepareMem(t *testing.T, ns string) Service {
	s := NewMemService()

	if err := s.Teardown(ns); err != nil {
		t.Fatal(err)
	}

	return s
}
