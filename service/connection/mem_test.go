package connection

import "testing"

func TestMemPut(t *testing.T) {
	testServicePut(t, prepareMem)
}

func TestMemPutInvalid(t *testing.T) {
	testServicePutInvalid(t, prepareMem)
}

func TestMemQuery(t *testing.T) {
	testServiceQuery(t, prepareMem)
}

func prepareMem(t *testing.T, ns string) Service {
	return NewMemService()
}
