package controller

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
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

	list, err := c.List(app, post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list.Comments), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, comment := range testCommentSet(owner.ID, post) {
		_, err = c.objects.Put(app.Namespace(), comment)
		if err != nil {
			t.Fatal(err)
		}
	}

	list, err = c.List(app, post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list.Comments), 5; have != want {
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

func TestCommentControllerExternalCreate(t *testing.T) {
	var (
		app, owner, c = testSetupCommentController(t)
		externalID    = "random-123_externalId"
	)

	created, err := c.ExternalCreate(app, owner, externalID, "Do Like.")
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

func TestCommentControllerExternalDelete(t *testing.T) {
	var (
		app, owner, c = testSetupCommentController(t)
		externalID    = "random-321_externalDeleteID"
	)

	created, err := c.objects.Put(
		app.Namespace(),
		testExternalComment(owner.ID, externalID),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ExternalDelete(app, owner, externalID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ExternalRetrieve(app, owner, externalID, created.ID)
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

	err = c.ExternalDelete(app, owner, externalID, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommentControllerExternallList(t *testing.T) {
	var (
		app, owner, c = testSetupCommentController(t)
		externalID    = "random-213_listID"
	)

	list, err := c.ExternalList(app, externalID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list.Comments), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, comment := range testExternalCommentSet(owner.ID, externalID) {
		_, err = c.objects.Put(app.Namespace(), comment)
		if err != nil {
			t.Fatal(err)
		}
	}

	list, err = c.ExternalList(app, externalID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(list.Comments), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestExternalCommentControllerRetrieve(t *testing.T) {
	var (
		app, owner, c = testSetupCommentController(t)
		externalID    = "external_retrieveID"
	)

	created, err := c.objects.Put(app.Namespace(), testExternalComment(owner.ID, externalID))
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.ExternalRetrieve(app, owner, externalID, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r, created; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %v, want %v", have, want)
	}

	_, err = c.ExternalRetrieve(app, owner, externalID, created.ID-1)
	if have, want := err, ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCommentExternalControllerUpdate(t *testing.T) {
	var (
		app, owner, c = testSetupCommentController(t)
		externalID    = "external-update-id"
	)

	_, err := c.ExternalUpdate(app, owner, externalID, 0, "Do not like.")
	if have, want := err, ErrNotFound; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	created, err := c.objects.Put(app.Namespace(), testExternalComment(owner.ID, externalID))
	if err != nil {
		t.Fatal(err)
	}

	updated, err := c.ExternalUpdate(app, owner, externalID, created.ID, "Do not like!")
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
		users   = user.NewNopService()
		user    = &v04_entity.ApplicationUser{
			ID: uint64(rand.Int63()),
		}
	)

	err := objects.Setup(app.Namespace())
	if err != nil {
		t.Fatal(err)
	}

	return app, user, NewCommentController(objects, users)
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

func testExternalComment(ownerID uint64, externalID string) *object.Object {
	return &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment("content", "Do like."),
		},
		ExternalID: externalID,
		OwnerID:    ownerID,
		Owned:      true,
		Type:       typeComment,
		Visibility: object.VisibilityPublic,
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

func testExternalCommentSet(ownerID uint64, externalID string) []*object.Object {
	return []*object.Object{
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ExternalID: externalID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: object.VisibilityPublic,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ExternalID: externalID,
			OwnerID:    ownerID + 1,
			Owned:      true,
			Type:       typeComment,
			Visibility: object.VisibilityPublic,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ExternalID: externalID,
			OwnerID:    ownerID - 1,
			Owned:      true,
			Type:       typeComment,
			Visibility: object.VisibilityPublic,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ExternalID: externalID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: object.VisibilityPublic,
		},
		{
			Attachments: []object.Attachment{
				object.NewTextAttachment("content", "Do like."),
			},
			ExternalID: externalID,
			OwnerID:    ownerID,
			Owned:      true,
			Type:       typeComment,
			Visibility: object.VisibilityPublic,
		},
	}
}