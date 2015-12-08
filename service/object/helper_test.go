package object

import "testing"

type prepareFunc func(string, *testing.T) Service

var testArticle = &Object{
	OwnerID:    555,
	Type:       "article",
	Visibility: VisibilityGlobal,
}

var testInvalid = &Object{
	Attachments: []Attachment{
		{
			Content: "foo barbaz",
			Name:    "summary",
			Type:    "invalid",
		},
	},
	Type:       "test",
	Visibility: VisibilityPrivate,
}

var testPost = &Object{
	Attachments: []Attachment{
		NewTextAttachment("intro", "Cupcake ipsum dolor sit amet."),
		NewURLAttachment("teaser", "http://bit.ly/1Jp8bMP"),
	},
	OwnerID:    123,
	Tags:       []string{"guide", "diy"},
	Type:       "post",
	Visibility: VisibilityConnection,
}

var testRecipe = &Object{
	Attachments: []Attachment{
		NewTextAttachment("yum", "Cupcake ipsum dolor sit amet."),
	},
	OwnerID:    321,
	Tags:       []string{"low-carb", "cold"},
	Type:       "recipe",
	Visibility: VisibilityConnection,
}

func testCreateSet(objectID uint64) []*Object {
	set := []*Object{}

	for i := 0; i < 5; i++ {
		set = append(set, &Object{
			OwnerID:    1,
			Type:       "article",
			Visibility: VisibilityConnection,
		})
	}

	for i := 0; i < 5; i++ {
		set = append(set, &Object{
			OwnerID:    1,
			Type:       "review",
			Visibility: VisibilityPublic,
		})
	}

	for i := 0; i < 5; i++ {
		set = append(set, &Object{
			OwnerID:    2,
			ObjectID:   objectID,
			Type:       "comment",
			Visibility: VisibilityGlobal,
		})
	}

	for i := 0; i < 13; i++ {
		set = append(set, &Object{
			OwnerID:    4,
			ObjectID:   objectID,
			Owned:      true,
			Type:       "tg_comment",
			Visibility: VisibilityConnection,
		})
	}

	return set
}

func testServiceQuery(t *testing.T, p prepareFunc) {
	var (
		namespace  = "service_query"
		service    = p(namespace, t)
		testObject = *testArticle

		owned bool
	)

	article, err := service.Put(namespace, &testObject)
	if err != nil {
		t.Fatal(err)
	}

	for _, o := range testCreateSet(article.ID) {
		_, err = service.Put(namespace, o)
		if err != nil {
			t.Fatal(err)
		}
	}

	os, err := service.Query(namespace, QueryOptions{
		OwnerIDs: []uint64{
			1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	os, err = service.Query(namespace, QueryOptions{
		ObjectIDs: []uint64{
			article.ID,
		},
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	os, err = service.Query(namespace, QueryOptions{
		Visibilities: []Visibility{
			VisibilityPublic,
			VisibilityGlobal,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 11; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	owned = true

	os, err = service.Query(namespace, QueryOptions{
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	os, err = service.Query(namespace, QueryOptions{
		Owned: &owned,
		Types: []string{
			"tg_comment",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 13; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = service.Query("invalid", QueryOptions{})
	if have, want := err, ErrNamespaceNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
