package connection

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestListFromIDs(t *testing.T) {
	var (
		ids = []uint64{
			uint64(rand.Int63()),
			uint64(rand.Int63()),
			uint64(rand.Int63()),
		}
		cons = List{}
	)

	for _, id := range ids {
		cons = append(cons, &Connection{FromID: id})
	}

	if have, want := cons.FromIDs(), ids; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestListToIDs(t *testing.T) {
	var (
		ids = []uint64{
			uint64(rand.Int63()),
			uint64(rand.Int63()),
			uint64(rand.Int63()),
		}
		cons = List{}
	)

	for _, id := range ids {
		cons = append(cons, &Connection{ToID: id})
	}

	if have, want := cons.ToIDs(), ids; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}
