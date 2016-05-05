package object

import (
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/asaskevich/govalidator"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/platform/service"
)

// Attachment variants available for Objects.
const (
	AttachmentTypeText = "text"
	AttachmentTypeURL  = "url"
)

// DefaultLanguage is used when no lang is provided for object content.
const DefaultLanguage = "en"

// Visibility variants available for Objects.
const (
	VisibilityPrivate Visibility = (iota + 1) * 10
	VisibilityConnection
	VisibilityPublic
	VisibilityGlobal
)

// Attachment is typed media which belongs to an Object.
type Attachment struct {
	Contents Contents `json:"contents"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
}

// Validate returns an error if a Attachment constraint is not full-filled.
func (a Attachment) Validate() error {
	if a.Name == "" {
		return wrapError(ErrInvalidAttachment, "name must be set")
	}

	if a.Type == "" ||
		(a.Type != AttachmentTypeText && a.Type != AttachmentTypeURL) {
		return wrapError(ErrInvalidAttachment, "unsupported type '%s'", a.Type)
	}

	if a.Contents == nil || len(a.Contents) == 0 {
		return wrapError(ErrInvalidAttachment, "contents can't be empty")
	}

	for tag, content := range a.Contents {
		_, err := language.Parse(tag)
		if err != nil {
			return wrapError(
				ErrInvalidAttachment,
				"invalid language tag '%s'",
				tag,
			)
		}

		if content == "" {
			return wrapError(ErrInvalidAttachment, "content missing for '%s'", tag)
		}

		if a.Type == AttachmentTypeURL && !govalidator.IsURL(content) {
			return wrapError(ErrInvalidAttachment, "invalid url for '%s'", tag)
		}
	}

	return nil
}

// NewTextAttachment returns an Attachment of type Text.
func NewTextAttachment(name string, contents Contents) Attachment {
	return Attachment{
		Contents: contents,
		Name:     name,
		Type:     AttachmentTypeText,
	}
}

// NewURLAttachment returns an Attachment of type URL.
func NewURLAttachment(name string, contents Contents) Attachment {
	return Attachment{
		Contents: contents,
		Name:     name,
		Type:     AttachmentTypeURL,
	}
}

// Contents is the mapping of content to locale.
type Contents map[string]string

// Validate performs semantic checks on the localisation fields.
func (c Contents) Validate() error {
	return nil
}

// List is an Object collection.
type List []*Object

// OwnerIDs returns all user ids of the associated object owners.
func (os List) OwnerIDs() []uint64 {
	ids := []uint64{}

	for _, o := range os {
		ids = append(ids, o.OwnerID)
	}

	return ids
}

// Map is an Object collection indexed by id.
type Map map[uint64]*Object

// Object is a generic building block to express different domains like Posts,
// Albums with their dependend objects.
type Object struct {
	Attachments []Attachment `json:"attachments"`
	CreatedAt   time.Time    `json:"created_at"`
	Deleted     bool         `json:"deleted"`
	ExternalID  string       `json:"external_id"`
	ID          uint64       `json:"id"`
	Latitude    float64      `json:"latitude"`
	Location    string       `json:"location"`
	Longitude   float64      `json:"longitude"`
	ObjectID    uint64       `json:"object_id"`
	Owned       bool         `json:"owned"`
	OwnerID     uint64       `json:"owner_id"`
	Private     *Private     `json:"private,omitempty"`
	Tags        []string     `json:"tags"`
	Type        string       `json:"type"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Visibility  Visibility   `json:"visibility"`
}

// Validate returns an error if a constraint on the Object is not full-filled.
func (o *Object) Validate() error {
	if len(o.Attachments) > 5 {
		return wrapError(ErrInvalidObject, "too many attachments")
	}

	for _, a := range o.Attachments {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	if o.OwnerID == 0 {
		return wrapError(ErrInvalidObject, "missing owner")
	}

	if len(o.Tags) > 5 {
		return wrapError(ErrInvalidObject, "too many tags")
	}

	if o.Type == "" {
		return wrapError(ErrInvalidObject, "missing type")
	}

	if o.Visibility < 10 || o.Visibility > 40 {
		return wrapError(ErrInvalidObject, "unsupported visibility")
	}

	return nil
}

// Private is the bucket for protected fields on an Object.
type Private struct {
	Visible bool `json:"visible"`
}

// QueryOptions are passed to narrow down query for objects.
type QueryOptions struct {
	Deleted      bool
	ExternalIDs  []string
	ID           *uint64
	ObjectIDs    []uint64
	OwnerIDs     []uint64
	Owned        *bool
	Types        []string
	Visibilities []Visibility
}

// Service for object interactions.
type Service interface {
	metrics.BucketByDay
	service.Lifecycle

	Count(namespace string, opts QueryOptions) (int, error)
	Put(namespace string, object *Object) (*Object, error)
	Query(namespace string, opts QueryOptions) (List, error)
	Remove(namespace string, id uint64) error
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// Visibility determines the visibility of Objects when consumed.
type Visibility uint8

func flakeNamespace(ns string) string {
	return fmt.Sprintf("%s_%s", ns, "objects")
}
