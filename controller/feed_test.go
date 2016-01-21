package controller

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/tapglue/multiverse/service/event"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestAffiliation(t *testing.T) {
	var (
		from = uint64(123)
		to   = uint64(321)
		a    = affiliations{
			&v04_entity.Connection{
				UserFromID: from,
				UserToID:   to,
				Type:       v04_entity.ConnectionTypeFollow,
			}: &v04_entity.ApplicationUser{
				ID: to,
			},
			&v04_entity.Connection{
				UserFromID: to,
				UserToID:   from,
				Type:       v04_entity.ConnectionTypeFollow,
			}: &v04_entity.ApplicationUser{
				ID: from,
			},
			&v04_entity.Connection{
				UserFromID: from,
				UserToID:   to,
				Type:       v04_entity.ConnectionTypeFriend,
			}: &v04_entity.ApplicationUser{
				ID: from,
			},
		}
	)

	if have, want := len(a.connections()), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.followers(from)), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.followings(from)), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.filterFollowers(from)), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.friends(from)), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.userIDs()), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := len(a.users()), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCollect(t *testing.T) {
	es, err := collect(
		testSourceLen(2),
		testSourceLen(7),
		testSourceLen(4),
	)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCollectError(t *testing.T) {
	_, err := collect(testSourceError)
	if err == nil {
		t.Error("want collect to error")
	}
}

func TestConditionDuplicate(t *testing.T) {
	es, err := testSourceLen(10)()
	if err != nil {
		t.Fatal(err)
	}

	es = append(es, &v04_entity.Event{
		ID: 5,
	})

	es = filter(es, conditionDuplicate())

	if have, want := len(es), 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestConditionObjectMissing(t *testing.T) {
	es, err := testSourceLen(10)()
	if err != nil {
		t.Fatal(err)
	}

	pm := PostMap{
		1: {},
		6: {},
	}

	es = filter(es, conditionPostMissing(pm))

	if have, want := len(es), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestFilter(t *testing.T) {
	es, err := testSourceLen(10)()
	if err != nil {
		t.Fatal(err)
	}

	es = filter(es, testConditionEven)

	if have, want := len(es), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestSourceConnection(t *testing.T) {
	var (
		from = uint64(123)
		to   = uint64(321)
		now  = time.Now()
		cs   = []*v04_entity.Connection{
			{
				State:      v04_entity.ConnectionStateConfirmed,
				Type:       v04_entity.ConnectionTypeFriend,
				UserFromID: from,
				UserToID:   to,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
			{
				State:      v04_entity.ConnectionStatePending,
				Type:       v04_entity.ConnectionTypeFollow,
				UserFromID: from,
				UserToID:   to,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
			{
				State:      v04_entity.ConnectionStateRejected,
				Type:       v04_entity.ConnectionTypeFollow,
				UserFromID: from,
				UserToID:   to,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
			{
				State:      v04_entity.ConnectionStateConfirmed,
				Type:       v04_entity.ConnectionTypeFollow,
				UserFromID: from,
				UserToID:   to,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
			{
				State:      v04_entity.ConnectionStateConfirmed,
				Type:       v04_entity.ConnectionTypeFollow,
				UserFromID: to,
				UserToID:   from,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
		}
	)

	es, err := sourceConnection(cs)()
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(es), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := es[0].Type, v04_entity.TypeEventFriend; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := es[1].Type, v04_entity.TypeEventFollow; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testConditionEven(idx int, event *v04_entity.Event) bool {
	return idx%2 == 0
}

func testSourceLen(n int) source {
	return func() (event.List, error) {
		es := event.List{}

		for i := 0; i < n; i++ {
			es = append(es, &v04_entity.Event{
				ID:       uint64(i + 1),
				ObjectID: uint64(i),
				Target: &v04_entity.Object{
					ID:   strconv.FormatUint(uint64(i+1), 10),
					Type: v04_entity.TypeTargetUser,
				},
				UserID: uint64(i + 1),
			})
		}

		return es, nil
	}
}

func testSourceError() (event.List, error) {
	return nil, fmt.Errorf("something went wrong")
}
