package controller

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestPostControllerCreate(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = &Post{
			Object: &object.Object{
				Attachments: []object.Attachment{
					object.NewTextAttachment("body", "Test body."),
				},
				Tags: []string{
					"review",
				},
				Visibility: object.VisibilityPublic,
			},
		}
	)

	created, err := c.Create(app, post, owner)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &created.ID,
		Owned: &defaultOwned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rs), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	if have, want := rs[0], created.Object; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerDelete(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = testPost(owner.ID)
	)

	created, err := c.objects.Put(app.Namespace(), post.Object)
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
		Owned:   &defaultOwned,
		Types: []string{
			typePost,
		},
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

func TestPostControllerListAll(t *testing.T) {
	app, owner, c := testSetupPostController(t)

	ps, err := c.ListAll(app)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ps), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, post := range testPostSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), post)
		if err != nil {
			t.Fatal(err)
		}
	}

	ps, err = c.ListAll(app)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ps), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerListUser(t *testing.T) {
	app, owner, c := testSetupPostController(t)

	ps, err := c.ListUser(app, owner.ID, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ps), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, post := range testPostSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), post)
		if err != nil {
			t.Fatal(err)
		}
	}

	ps, err = c.ListUser(app, owner.ID, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ps), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerRetrieve(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = testPost(owner.ID)
	)

	created, err := c.objects.Put(app.Namespace(), post.Object)
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Retrieve(app, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r.Object, created; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %v, want %v", have, want)
	}

	_, err = c.Retrieve(app, created.ID-1)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerUpdate(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = testPost(owner.ID)
	)

	created, err := c.objects.Put(app.Namespace(), post.Object)
	if err != nil {
		t.Fatal(err)
	}

	created.OwnerID = 0

	_, err = c.Update(app, owner, created.ID, &Post{Object: created})
	if err != nil {
		t.Fatal(err)
	}

	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &created.ID,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ps), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	updated := ps[0]

	if have, want := updated.OwnerID, post.OwnerID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := updated.Visibility, post.Visibility; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerUpdateMissing(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = testPost(owner.ID)
	)

	_, err := c.Update(app, owner, post.ID, post)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testSetupPostController(
	t *testing.T,
) (*v04_entity.Application, *v04_entity.ApplicationUser, *PostController) {
	var (
		app = &v04_entity.Application{
			ID:    rand.Int63(),
			OrgID: rand.Int63(),
		}
		objects = object.NewMemService()
		user    = &v04_entity.ApplicationUser{
			ID: uint64(rand.Int63()),
		}
	)

	err := objects.Setup(app.Namespace())
	if err != nil {
		t.Fatal(err)
	}

	return app, user, NewPostController(connection.NewNopService(), objects)
}

func testPost(ownerID uint64) *Post {
	return &Post{
		Object: &object.Object{
			Attachments: []object.Attachment{
				object.NewTextAttachment("body", "Test body."),
			},
			OwnerID: ownerID,
			Owned:   true,
			Tags: []string{
				"review",
			},
			Type:       typePost,
			Visibility: object.VisibilityPublic,
		},
	}
}

func testPostSet(ownerID uint64) []*object.Object {
	return []*object.Object{
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typePost,
			Visibility: object.VisibilityConnection,
		},
		{
			OwnerID:    ownerID + 1,
			Owned:      true,
			Type:       typePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID - 1,
			Owned:      true,
			Type:       typePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typePost,
			Visibility: object.VisibilityPrivate,
		},
	}
}
