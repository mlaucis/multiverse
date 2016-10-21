// +build integration

package member

import (
	"flag"
	"fmt"
	"math/rand"
	"os/user"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var pgTestURL string

func TestPostgresPut(t *testing.T) {
	var (
		enabled   = true
		namespace = "service_put"
		service   = preparePostgres(t, namespace)
		member    = testMember()
	)

	created, err := service.Put(namespace, member)
	if err != nil {
		t.Fatal(err)
	}

	list, err := service.Query(namespace, QueryOpts{
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

	list, err = service.Query(namespace, QueryOpts{
		Enabled: &created.Enabled,
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
	if have, want := list[0], updated; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = service.Put(namespace, &Member{})
	if have, want := err, ErrInvalidMember; !IsInvalidMember(err) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostgresQuery(t *testing.T) {
	t.Fatal("TestQuery not implemented")
}

func testMember() *Member {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Member{
		Email: fmt.Sprintf(
			"member%d@tapglue.test",
			r.Int63(),
		),
		Enabled: true,
	}
}

func preparePostgres(t *testing.T, namespace string) Service {
	db, err := sqlx.Connect("postgres", pgTestURL)
	if err != nil {
		t.Fatal(err)
	}

	s := PostgresService(db)

	if err := s.Teardown(namespace); err != nil {
		t.Fatal(err)
	}

	return s
}

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	d := fmt.Sprintf(
		"postgres://%s@127.0.0.1:5432/tapglue_test?sslmode=disable&connect_timeout=5",
		u.Username,
	)

	url := flag.String("postgres.url", d, "Postgreas integration connection URL")
	flag.Parse()

	pgTestURL = *url
}
