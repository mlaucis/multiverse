// build +integration

package event

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/user"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	day  = 24 * time.Hour
	week = 7 * day

	pgURL string
)

type testEvent struct {
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TestActiveUserIDs(t *testing.T) {
	var (
		ns     = "active_user"
		insert = wrapNamespace(`INSERT INTO %s.events(json_data) VALUES($1)`, ns)
		s, db  = preparePostgres(ns, t)
	)

	testSet := map[uint64]map[time.Duration]int{
		123: map[time.Duration]int{
			time.Hour: 13,
			day:       3,
		},
		321: map[time.Duration]int{
			time.Hour: 3,
			day:       15,
		},
		456: map[time.Duration]int{
			day:  3,
			week: 21,
		},
	}

	for id, periods := range testSet {
		for d, amount := range periods {
			for i := 0; i < amount; i++ {
				data, err := json.Marshal(&testEvent{
					UserID:    id,
					UpdatedAt: time.Now().Add(-d),
				})
				if err != nil {
					t.Fatal(err)
				}

				_, err = db.Exec(insert, data)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	ids, err := s.ActiveUserIDs(ns, ByDay)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ids), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ids, []uint64{123, 321}; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	ids, err = s.ActiveUserIDs(ns, ByWeek)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ids), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ids, []uint64{321, 123, 456}; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	ids, err = s.ActiveUserIDs(ns, ByMonth)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ids), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ids, []uint64{456, 321, 123}; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func preparePostgres(namespace string, t *testing.T) (Service, *sqlx.DB) {
	db, err := sqlx.Connect("postgres", pgURL)
	if err != nil {
		t.Fatal(err)
	}

	s := NewPostgresService(db)

	err = s.Teardown(namespace)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Setup(namespace)
	if err != nil {
		t.Fatal(err)
	}

	return s, db
}

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	d := fmt.Sprintf(
		"postgres://%s@127.0.0.1:5432/tapglue_test?sslmode=disable&connect_timeout=5",
		user.Username,
	)

	url := flag.String("postgres.url", d, "Postgres connection URL")
	flag.Parse()

	pgURL = *url
}
