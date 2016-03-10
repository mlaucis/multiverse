package event

import (
	"reflect"
	"testing"
)

type prepareFunc func(namespace string, t *testing.T) Service

func testEvent() *Event {
	return &Event{
		Enabled: true,
	}
}

func testList(objectID uint64, externalID string) List {
	es := List{}

	for i := 0; i < 5; i++ {
		es = append(es, &Event{
			Enabled:    true,
			Type:       "bookmark",
			UserID:     1,
			Visibility: VisibilityConnection,
		})
	}

	for i := 0; i < 5; i++ {
		es = append(es, &Event{
			Enabled:    true,
			ObjectID:   objectID,
			Type:       "share",
			UserID:     2,
			Visibility: VisibilityConnection,
		})
	}

	for i := 0; i < 5; i++ {
		es = append(es, &Event{
			Enabled:    true,
			Type:       "play",
			UserID:     2,
			Visibility: VisibilityPublic,
		})
	}

	for i := 0; i < 5; i++ {
		es = append(es, &Event{
			Enabled:    true,
			Type:       "vote",
			UserID:     2,
			Visibility: VisibilityPublic,
		})
	}

	for i := 0; i < 5; i++ {
		es = append(es, &Event{
			Enabled:    true,
			ObjectID:   objectID,
			Owned:      true,
			Type:       "tg_share",
			UserID:     4,
			Visibility: VisibilityConnection,
		})
	}

	for i := 0; i < 6; i++ {
		es = append(es, &Event{
			Enabled:    true,
			ObjectID:   objectID,
			Owned:      true,
			Type:       "tg_like",
			UserID:     4,
			Visibility: VisibilityConnection,
		})
	}

	for i := 0; i < 11; i++ {
		es = append(es, &Event{
			Enabled: true,
			Object: &Object{
				ID:   externalID,
				Type: "restaurant",
			},
			Type:       "checkin",
			UserID:     4,
			Visibility: VisibilityConnection,
		})
	}

	return es
}

func testServicePut(p prepareFunc, t *testing.T) {
	var (
		event     = testEvent()
		namespace = "service_put"
		service   = p(namespace, t)
		enabled   = true
	)

	created, err := service.Put(namespace, event)
	if err != nil {
		t.Fatal(err)
	}

	list, err := service.Query(namespace, QueryOptions{
		Enabled: &enabled,
		IDs: []uint64{
			created.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := list[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	created.Enabled = false

	updated, err := service.Put(namespace, created)
	if err != nil {
		t.Fatal(err)
	}

	list, err = service.Query(namespace, QueryOptions{
		Enabled: &created.Enabled,
		IDs: []uint64{
			created.ID,
		},
	})

	if have, want := len(list), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := list[0], updated; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = service.Put("invalid", event)
	if have, want := err, ErrNamespaceNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testServiceQuery(p prepareFunc, t *testing.T) {
	var (
		namespace         = "service_query"
		service           = p(namespace, t)
		enabled           = true
		externalID        = "external-id-123"
		objectID   uint64 = 321
		owned             = false
	)

	for _, e := range testList(objectID, externalID) {
		_, err := service.Put(namespace, e)
		if err != nil {
			t.Fatal(err)
		}
	}

	es, err := service.Query(namespace, QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 42; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		ObjectIDs: []uint64{
			objectID,
		},
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		UserIDs: []uint64{
			1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		Visibilities: []Visibility{
			VisibilityPublic,
			VisibilityGlobal,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	owned = true

	es, err = service.Query(namespace, QueryOptions{
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 11; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		Owned: &owned,
		Types: []string{
			"tg_like",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 6; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		ExternalObjectIDs: []string{
			externalID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 11; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	es, err = service.Query(namespace, QueryOptions{
		ExternalObjectTypes: []string{
			"restaurant",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 11; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
