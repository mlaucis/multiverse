// +build integration

package object

import (
	"fmt"
	"os/user"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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

	err = service.Remove("invalid_namespace", 123)
	if have, want := err, ErrNamespaceNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func preparePostgres(namespace string, t *testing.T) Service {
	user, err := user.Current()
	if err != nil {
		t.Fatal(t)
	}

	url := fmt.Sprintf(
		"postgres://%s@127.0.0.1:5432/tapglue_test?sslmode=disable&connect_timeout=5",
		user.Username,
	)

	db, err := sqlx.Connect("postgres", url)
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
