// +build integration

package object

import (
	"flag"
	"fmt"
	"os/user"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var pgURL string

func TestPostgresServiceCount(t *testing.T) {
	testServiceCount(t, preparePostgres)
}

func TestPostgresServicePut(t *testing.T) {
	var (
		namespace = "service_put"
		service   = preparePostgres(namespace, t)
		post      = *testPost
	)

	created, err := service.Put(namespace, &post)
	if err != nil {
		t.Fatal(err)
	}

	list, err := service.Query(namespace, QueryOptions{
		ID: &created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := list[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v, want %#v", have, want)
	}

	created.Deleted = true

	updated, err := service.Put(namespace, created)
	if err != nil {
		t.Fatal(err)
	}

	list, err = service.Query(namespace, QueryOptions{
		Deleted: true,
		ID:      &created.ID,
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
}

func TestPostgresServicePutInvalid(t *testing.T) {
	var (
		namespace = "service_put_invalid"
		service   = preparePostgres(namespace, t)
		invalid   = *testInvalid
	)

	_, err := service.Put(namespace, &invalid)
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestPostgresServiceQuery(t *testing.T) {
	testServiceQuery(t, preparePostgres)
}

func TestPostgresServiceRemove(t *testing.T) {
	var (
		namespace = "service_remove"
		service   = preparePostgres(namespace, t)
		recipe    = *testRecipe
	)

	created, err := service.Put(namespace, &recipe)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Remove(namespace, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	list, err := service.Query(namespace, QueryOptions{
		ID: &created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	err = service.Remove(namespace, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func preparePostgres(namespace string, t *testing.T) Service {
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

	return s
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
