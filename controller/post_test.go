package controller

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestPostControllerCreate(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = &Post{
			Object: &object.Object{
				Attachments: []object.Attachment{
					object.NewTextAttachment("body", object.Contents{
						"en": "Test body.",
					}),
				},
				Tags: []string{
					"review",
				},
				Visibility: object.VisibilityPublic,
			},
		}
	)

	created, err := c.Create(
		app,
		Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		},
		post,
	)
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

func TestPostCreateConstrainVisibility(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = &Post{
			Object: &object.Object{
				Visibility: object.VisibilityGlobal,
			},
		}
	)

	_, err := c.Create(
		app,
		Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		},
		post,
	)

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
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

	err = c.Delete(app, owner.ID+1, created.ID)
	if have, want := err, ErrUnauthorized; !IsUnauthorized(err) {
		t.Errorf("have %v, want %v", have, want)
	}

	err = c.Delete(app, owner.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Retrieve(app, owner.ID, created.ID)
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		Deleted: true,
		ID:      &created.ID,
		Owned:   &defaultOwned,
		Types: []string{
			TypePost,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	err = c.Delete(app, owner.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostControllerListAll(t *testing.T) {
	app, owner, c := testSetupPostController(t)

	feed, err := c.ListAll(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(feed.Posts), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, post := range testPostSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), post)
		if err != nil {
			t.Fatal(err)
		}
	}

	feed, err = c.ListAll(app, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(feed.Posts), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerListUser(t *testing.T) {
	app, owner, c := testSetupPostController(t)

	feed, err := c.ListUser(app, owner.ID, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(feed.Posts), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, post := range testPostSet(owner.ID) {
		_, err = c.objects.Put(app.Namespace(), post)
		if err != nil {
			t.Fatal(err)
		}
	}

	feed, err = c.ListUser(app, owner.ID, owner.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(feed.Posts), 3; have != want {
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

	r, err := c.Retrieve(app, owner.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r.Object, created; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %v, want %v", have, want)
	}

	_, err = c.Retrieve(app, owner.ID, created.ID-1)
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

	_, err = c.Update(
		app,
		Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		},
		created.ID,
		&Post{Object: created},
	)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &created.ID,
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
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

func TestPostUpdateConstrainVisibility(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		origin        = Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		}
		post = testPost(owner.ID)
	)

	created, err := c.Create(
		app,
		origin,
		post,
	)
	if err != nil {
		t.Fatal(err)
	}

	created.Visibility = object.VisibilityGlobal

	_, err = c.Update(app, origin, created.ID, created)

	if have, want := err, ErrUnauthorized; !IsUnauthorized(have) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostControllerUpdateMissing(t *testing.T) {
	var (
		app, owner, c = testSetupPostController(t)
		post          = testPost(owner.ID)
	)

	_, err := c.Update(
		app,
		Origin{
			Integration: IntegrationApplication,
			UserID:      owner.ID,
		},
		post.ID,
		post,
	)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testSetupPostController(
	t *testing.T,
) (*v04_entity.Application, *user.User, *PostController) {
	var (
		app = &v04_entity.Application{
			ID:    rand.Int63(),
			OrgID: rand.Int63(),
		}
		events  = event.NewMemService()
		objects = object.NewMemService()
		users   = user.NewMemService()
		u       = &user.User{
			ID: uint64(rand.Int63()),
		}
	)

	err := events.Setup(app.Namespace())
	if err != nil {
		t.Fatal(err)
	}

	err = objects.Setup(app.Namespace())
	if err != nil {
		t.Fatal(err)
	}

	return app, u, NewPostController(
		connection.NewMemService(),
		events,
		objects,
		users,
	)
}

func testPost(ownerID uint64) *Post {
	return &Post{
		Object: &object.Object{
			Attachments: []object.Attachment{
				object.NewTextAttachment("body", object.Contents{
					"en": "Test body.",
				}),
			},
			OwnerID: ownerID,
			Owned:   true,
			Tags: []string{
				"review",
			},
			Type:       TypePost,
			Visibility: object.VisibilityPublic,
		},
	}
}

func testPostSet(ownerID uint64) []*object.Object {
	return []*object.Object{
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       TypePost,
			Visibility: object.VisibilityConnection,
		},
		{
			OwnerID:    ownerID + 1,
			Owned:      true,
			Type:       TypePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID - 1,
			Owned:      true,
			Type:       TypePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       TypePost,
			Visibility: object.VisibilityPublic,
		},
		{
			OwnerID:    ownerID,
			Owned:      true,
			Type:       TypePost,
			Visibility: object.VisibilityPrivate,
		},
	}
}
