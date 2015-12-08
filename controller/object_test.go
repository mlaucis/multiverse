package controller

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestObjectControllerCreate(t *testing.T) {
	var (
		app, owner, c = testSetupObjectController(t)
		recipe        = &object.Object{
			Type:       "recipe",
			Visibility: object.VisibilityPrivate,
		}
	)

	created, err := c.Create(app, recipe, owner)
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Retrieve(app, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r.OwnerID, owner.ID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerDelete(t *testing.T) {
	var (
		app, owner, c = testSetupObjectController(t)
		article       = &object.Object{
			OwnerID:    owner.ID,
			Type:       "article",
			Visibility: object.VisibilityPublic,
		}
	)

	created, err := c.objects.Put(app.Namespace(), article)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Delete(app, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Retrieve(app, created.ID)
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		Deleted: true,
		ID:      &created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	err = c.Delete(app, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestObjectControllerList(t *testing.T) {
	app, owner, c := testSetupObjectController(t)

	as, err := c.List(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(as), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, article := range testArticleSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), article)
		if err != nil {
			t.Fatal(err)
		}
	}

	as, err = c.List(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(as), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerListAll(t *testing.T) {
	app, owner, c := testSetupObjectController(t)

	as, err := c.ListAll(app)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(as), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, article := range testArticleSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), article)
		if err != nil {
			t.Fatal(err)
		}
	}

	as, err = c.ListAll(app)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(as), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerListConnections(t *testing.T) {
	app, owner, c := testSetupObjectController(t)

	as, err := c.ListConnections(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(as), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, article := range testArticleSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), article)
		if err != nil {
			t.Fatal(err)
		}
	}

	os, err := c.ListConnections(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	// FIXME(xla): Populate connections service with proper IDs.
	if have, want := len(os), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerRetrieve(t *testing.T) {
	var (
		app, owner, c = testSetupObjectController(t)
		review        = &object.Object{
			OwnerID:    owner.ID,
			Type:       "review",
			Visibility: object.VisibilityPublic,
		}
	)

	created, err := c.objects.Put(app.Namespace(), review)
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Retrieve(app, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r, created; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %v, want %v", have, want)
	}

	_, err = c.Retrieve(app, created.ID-1)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerUpdate(t *testing.T) {
	var (
		app, _, c = testSetupObjectController(t)
		post      = &object.Object{
			OwnerID:    123,
			Type:       "post",
			Visibility: object.VisibilityGlobal,
		}
		ns = app.Namespace()
	)

	created, err := c.objects.Put(ns, post)
	if err != nil {
		t.Fatal(err)
	}

	created.OwnerID = 0

	_, err = c.Update(app, created.ID, created)
	if err != nil {
		t.Fatal(err)
	}

	os, err := c.objects.Query(ns, object.QueryOptions{
		ID: &created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	updated := os[0]

	if have, want := updated.OwnerID, post.OwnerID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := updated.Type, post.Type; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := updated.Visibility, post.Visibility; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestObjectControllerUpdateMissing(t *testing.T) {
	var (
		app, _, c = testSetupObjectController(t)
		post      = &object.Object{
			ID:         1010,
			OwnerID:    321,
			Type:       "post",
			Visibility: object.VisibilityGlobal,
		}
	)

	_, err := c.Update(app, post.ID, post)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testSetupObjectController(
	t *testing.T,
) (*v04_entity.Application, *v04_entity.ApplicationUser, *ObjectController) {
	var (
		objects = object.NewMemService()
		app     = &v04_entity.Application{
			ID:    rand.Int63(),
			OrgID: rand.Int63(),
		}
		user = &v04_entity.ApplicationUser{
			ID: uint64(rand.Int63()),
		}
	)

	err := objects.Setup(app.Namespace())
	if err != nil {
		t.Fatal(err)
	}

	return app, user, NewObjectController(connection.NewNopService(), objects)
}

func testArticleSet(id uint64) []*object.Object {
	return []*object.Object{
		{
			OwnerID:    id,
			Type:       "article",
			Visibility: object.VisibilityConnection,
		},
		{
			OwnerID:    id + 1,
			Type:       "article",
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    id - 1,
			Type:       "article",
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    id,
			Type:       "article",
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    id,
			Type:       "article",
			Visibility: object.VisibilityPrivate,
		},
	}
}
