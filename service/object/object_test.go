package object

import "testing"

func TestAttachmentValidate(t *testing.T) {
	for _, a := range []Attachment{
		// Missing Content
		{
			Content: "",
			Name:    "attach1",
			Type:    AttachmentTypeText,
		},
		// Missing Name
		{
			Content: "Lorem ipsum.",
			Name:    "",
			Type:    AttachmentTypeText,
		},
		// Missing Type
		{
			Content: "Lorem ipsum.",
			Name:    "teaser",
			Type:    "",
		},
		// Unspported Type
		{
			Content: "Lorem ipsum.",
			Name:    "teaser",
			Type:    "teaser",
		},
		// Invalid URL
		{
			Content: "http://bit.ly^fake",
			Name:    "attach2",
			Type:    AttachmentTypeURL,
		},
	} {
		if err := a.Validate(); err == nil {
			t.Errorf("expected error: %v", a)
		}
	}
}

func TestObjectValidate(t *testing.T) {
	for _, o := range []*Object{
		// Too many Attachments
		{
			Attachments: []Attachment{
				{},
				{},
				{},
				{},
				{},
				{},
			},
		},
		// Too many Tags
		{
			Tags: []string{
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
				"tag",
			},
			Type:       "post",
			Visibility: VisibilityConnection,
		},
		// Missing Type
		{
			Visibility: VisibilityConnection,
		},
		// Invalid Visibility
		{
			Type:       "recipe",
			Visibility: 60,
		},
	} {
		if err := o.Validate(); err == nil {
			t.Errorf("expected error: %v", o)
		}
	}
}
