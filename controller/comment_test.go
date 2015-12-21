package controller

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

func TestCommentControllerCreate(t *testing.T) {
	app, owner, c := testSetupCommentController(t)

	post, err := c.objects.Put(app.Namespace(), testPost(owner.ID).Object)
	if err != nil {
		t.Fatal(err)
	}

	created, err := c.Create(app, owner, post.ID, "Do Like.")
	if err != nil {
		t.Fatal(err)
	}

	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &created.ID,
		Owned: &defaultOwned,
		Types: []string{
			typeComment,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	if have, want := cs[0], created; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = c.Create(app, owner, 0, "Do not like.")
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
}

func TestCommentControllerDelete(t *testing.T) {
	app, owner, c := testSetupCommentController(t)

	post, err := c.objects.Put(app.Namespace(), testPost(owner.ID).Object)
	if err != nil {
		t.Fatal(err)
	}

	created, err := c.objects.Put(app.Namespace(), testComment(owner.ID, post))
	if err != nil {
		t.Fatal(err)
	}

	err = c.Delete(app, owner, post.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Retrieve(app, owner, post.ID, created.ID)
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		Deleted: true,
		ID:      &created.ID,
		Owned:   &defaultOwned,
		Types: []string{
			typeComment,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	err = c.Delete(app, owner, post.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentControllerList(t *testing.T) {
	app, owner, c := testSetupCommentController(t)

	post, err := c.objects.Put(app.Namespace(), testPost(owner.ID).Object)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := c.List(app, post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, comment := range testCommentSet(owner.ID, post) {
		_, err = c.objects.Put(app.Namespace(), comment)
		if err != nil {
			t.Fatal(err)
		}
	}

	cs, err = c.List(app, post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCommentControllerRetrieve(t *testing.T) {
	app, owner, c := testSetupCommentController(t)

	post, err := c.objects.Put(app.Namespace(), testPost(owner.ID).Object)
	if err != nil {
		t.Fatal(err)
	}

	created, err := c.objects.Put(app.Namespace(), testComment(owner.ID, post))
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Retrieve(app, owner, post.ID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r, created; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %v, want %v", have, want)
	}

	_, err = c.Retrieve(app, owner, post.ID, created.ID-1)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCommentControllerUpdate(t *testing.T) {
	app, owner, c := testSetupCommentController(t)

	post, err := c.objects.Put(app.Namespace(), testPost(owner.ID).Object)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Update(app, owner, post.ID, 0, "Do not like.")
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	created, err := c.objects.Put(app.Namespace(), testComment(owner.ID, post))
	if err != nil {
		t.Fatal(err)
	}

	updated, err := c.Update(app, owner, post.ID, created.ID, "Do not like!")
	if err != nil {
		t.Fatal(err)
	}

	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &created.ID,
		Owned: &defaultOwned,
		Types: []string{
			typeComment,
		},
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

func testSetupCommentController(
	t *testing.T,
) (*v04_entity.Application, *v04_entity.ApplicationUser, *CommentController) {
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

	return app, user, NewCommentController(objects)
}

func testComment(ownerID uint64, post *object.Object) *object.Object {
	return &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment("content", "Do like."),
		},
		ObjectID:   post.ID,
		OwnerID:    ownerID,
		Owned:      true,
		Type:       typeComment,
		Visibility: post.Visibility,
	}
}

func testCommentSet(ownerID uint64, post *object.Object) []*object.Object {
	return []*object.Object{
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ObjectID:   post.ID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: post.Visibility,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ObjectID:   post.ID,
			OwnerID:    ownerID + 1,
			Owned:      true,
			Type:       typeComment,
			Visibility: post.Visibility,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ObjectID:   post.ID,
			OwnerID:    ownerID - 1,
			Owned:      true,
			Type:       typeComment,
			Visibility: post.Visibility,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ObjectID:   post.ID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: post.Visibility,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ObjectID:   post.ID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: post.Visibility,
		},
	}
}
