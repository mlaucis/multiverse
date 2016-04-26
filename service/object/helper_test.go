package object

import "testing"

type prepareFunc func(string, *testing.T) Service

var testArticle = &Object{
	OwnerID:    555,
	Type:       "article",
	Visibility: VisibilityGlobal,
}

var testAttachmentText = NewTextAttachment("intro", Contents{
	"en": "Cupcake ipsum dolor sit amet.",
})

var testAttachmentURL = NewURLAttachment("teaser", Contents{
	"en": "http://bit.ly/1Jp8bMP",
})

var testInvalid = &Object{
	Attachments: []Attachment{
		{
			Contents: Contents{
				"en": "foo barbaz",
			},
			Name: "summary",
			Type: "invalid",
		},
	},
	Type:       "test",
	Visibility: VisibilityPrivate,
}

var testPost = &Object{
	Attachments: []Attachment{
		testAttachmentText,
		testAttachmentURL,
	},
	OwnerID:    123,
	Tags:       []string{"guide", "diy"},
	Type:       "post",
	Visibility: VisibilityConnection,
}

var testRecipe = &Object{
	Attachments: []Attachment{
		NewTextAttachment("yum", Contents{
			"en": "Cupcake ipsum dolor sit amet.",
		}),
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

	for i := 0; i < 7; i++ {
		set = append(set, &Object{
			ExternalID: "external-input-123",
			OwnerID:    5,
			Owned:      true,
			Type:       "tg_comment",
			Visibility: VisibilityConnection,
		})
	}

	return set
}

func testServiceCount(t *testing.T, p prepareFunc) {
	var (
		namespace  = "service_count"
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

	count, err := service.Count(namespace, QueryOptions{
		OwnerIDs: []uint64{
			1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 10; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		ObjectIDs: []uint64{
			article.ID,
		},
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		Visibilities: []Visibility{
			VisibilityPublic,
			VisibilityGlobal,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 11; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	owned = true

	count, err = service.Count(namespace, QueryOptions{
		Owned: &owned,
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 20; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		Owned: &owned,
		Types: []string{
			"tg_comment",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 20; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	count, err = service.Count(namespace, QueryOptions{
		ExternalIDs: []string{
			"external-input-123",
		},
		Owned: &owned,
		Types: []string{
			"tg_comment",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := count, 7; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
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

	if have, want := len(os), 20; have != want {
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

	if have, want := len(os), 20; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	os, err = service.Query(namespace, QueryOptions{
		ExternalIDs: []string{
			"external-input-123",
		},
		Owned: &owned,
		Types: []string{
			"tg_comment",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(os), 7; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	// FIXME(xla): Re-enable as soon as we return the error.
	// _, err = service.Query("invalid", QueryOptions{})
	// if have, want := err, ErrNamespaceNotFound; have != want {
	// 	t.Errorf("have %v, want %v", have, want)
	// }
}
