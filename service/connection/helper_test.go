package connection

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type prepareFunc func(t *testing.T, namespace string) Service

func testServiceCount(t *testing.T, p prepareFunc) {
	var (
		namespace = "service_count"
		service   = p(t, namespace)
		from      = uint64(rand.Int63())
		to        = uint64(rand.Int63())
		disabled  = false
	)

	for _, c := range testList(from, to) {
		_, err := service.Put(namespace, c)
		if err != nil {
			t.Fatal(err)
		}
	}

	c, err := service.Count(namespace, QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 25; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err = service.Count(namespace, QueryOptions{
		Enabled: &disabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err = service.Count(namespace, QueryOptions{
		FromIDs: []uint64{
			from,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 12; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	limit := 10

	c, err = service.Count(namespace, QueryOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err = service.Count(namespace, QueryOptions{
		States: []State{
			StateConfirmed,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 7; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err = service.Count(namespace, QueryOptions{
		ToIDs: []uint64{
			to,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err = service.Count(namespace, QueryOptions{
		Types: []Type{
			TypeFriend,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c, 18; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testList(from, to uint64) List {
	cs := List{}

	for i := 0; i < 7; i++ {
		cs = append(cs, &Connection{
			Enabled:   true,
			FromID:    from,
			State:     StateConfirmed,
			ToID:      uint64(rand.Int63()),
			Type:      TypeFollow,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	for i := 0; i < 5; i++ {
		cs = append(cs, &Connection{
			Enabled:   false,
			FromID:    from,
			State:     StatePending,
			ToID:      uint64(rand.Int63()),
			Type:      TypeFriend,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	for i := 0; i < 13; i++ {
		cs = append(cs, &Connection{
			Enabled:   true,
			FromID:    uint64(rand.Int63()),
			State:     StateRejected,
			ToID:      to,
			Type:      TypeFriend,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return cs
}

func testServicePut(t *testing.T, p prepareFunc) {
	var (
		namespace = "service_put"
		service   = p(t, namespace)
		con       = &Connection{
			Enabled: true,
			FromID:  uint64(rand.Int63()),
			ToID:    uint64(rand.Int63()),
			Type:    TypeFollow,
			State:   StatePending,
		}
	)

	created, err := service.Put(namespace, con)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := service.Query(namespace, QueryOptions{
		Enabled: &con.Enabled,
		FromIDs: []uint64{con.FromID},
		Types:   []Type{con.Type},
		ToIDs:   []uint64{con.ToID},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := cs[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	created.State = StateConfirmed

	updated, err := service.Put(namespace, created)
	if err != nil {
		t.Fatal(err)
	}

	cs, err = service.Query(namespace, QueryOptions{
		Enabled: &con.Enabled,
		FromIDs: []uint64{con.FromID},
		Types:   []Type{con.Type},
		ToIDs:   []uint64{con.ToID},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := cs[0], updated; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testServicePutInvalid(t *testing.T, p prepareFunc) {
	var (
		namespace = "service_put_invalid"
		service   = p(t, namespace)
	)

	// missing FromID
	_, err := service.Put(namespace, &Connection{})
	if !IsInvalidConnection(err) {
		t.Errorf("expected error: %s", ErrInvalidConnection)
	}

	// missing ToID
	_, err = service.Put(namespace, &Connection{
		FromID: uint64(rand.Int63()),
	})
	if !IsInvalidConnection(err) {
		t.Errorf("expected error: %s", ErrInvalidConnection)
	}

	// missing State
	_, err = service.Put(namespace, &Connection{
		FromID: uint64(rand.Int63()),
		ToID:   uint64(rand.Int63()),
	})
	if !IsInvalidConnection(err) {
		t.Errorf("expected error: %s", ErrInvalidConnection)
	}

	// missing Type
	_, err = service.Put(namespace, &Connection{
		FromID: uint64(rand.Int63()),
		ToID:   uint64(rand.Int63()),
		State:  StateConfirmed,
	})
	if !IsInvalidConnection(err) {
		t.Errorf("expected error: %s", ErrInvalidConnection)
	}
}

func testServiceQuery(t *testing.T, p prepareFunc) {
	var (
		namespace = "service_query"
		service   = p(t, namespace)
		from      = uint64(rand.Int63())
		to        = uint64(rand.Int63())
		disabled  = false
	)

	for _, c := range testList(from, to) {
		_, err := service.Put(namespace, c)
		if err != nil {
			t.Fatal(err)
		}
	}

	cs, err := service.Query(namespace, QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 25; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	cs, err = service.Query(namespace, QueryOptions{
		Enabled: &disabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	cs, err = service.Query(namespace, QueryOptions{
		FromIDs: []uint64{
			from,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 12; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	limit := 10

	cs, err = service.Query(namespace, QueryOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	cs, err = service.Query(namespace, QueryOptions{
		States: []State{
			StateConfirmed,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 7; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	cs, err = service.Query(namespace, QueryOptions{
		ToIDs: []uint64{
			to,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	cs, err = service.Query(namespace, QueryOptions{
		Types: []Type{
			TypeFriend,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 18; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
