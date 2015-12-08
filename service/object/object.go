package object

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

// Attachment variants available for Objects.
const (
	AttachmentTypeText = "text"
	AttachmentTypeURL  = "url"
)

// Visibility variants available for Objects.
const (
	VisibilityPrivate Visibility = (iota + 1) * 10
	VisibilityConnection
	VisibilityPublic
	VisibilityGlobal
)

// Attachment is typed media which belongs to an Object.
type Attachment struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

// Validate returns an error if a Attachment constraint is not full-filled.
func (a Attachment) Validate() error {
	if a.Content == "" {
		return ErrInvalidAttachment
	}

	if a.Name == "" {
		return ErrInvalidAttachment
	}

	if a.Type == "" ||
		(a.Type != AttachmentTypeText && a.Type != AttachmentTypeURL) {
		return ErrInvalidAttachment
	}

	if a.Type == AttachmentTypeURL {
		if !govalidator.IsURL(a.Content) {
			return ErrInvalidAttachment
		}
	}

	return nil
}

// NewTextAttachment returns an Attachment of type Text.
func NewTextAttachment(name, content string) Attachment {
	return Attachment{
		Content: content,
		Name:    name,
		Type:    AttachmentTypeText,
	}
}

// NewURLAttachment returns an Attachment of type URL.
func NewURLAttachment(name, content string) Attachment {
	return Attachment{
		Content: content,
		Name:    name,
		Type:    AttachmentTypeURL,
	}
}

// Object is a generic building block to express different domains like Posts,
// Albums with their dependend objects.
type Object struct {
	Attachments []Attachment `json:"attachments"`
	CreatedAt   time.Time    `json:"created_at"`
	Deleted     bool         `json:"deleted"`
	ID          uint64       `json:"id"`
	Latitude    float64      `json:"latitude"`
	Location    string       `json:"location"`
	Longitude   float64      `json:"longitude"`
	ObjectID    uint64       `json:"object_id"`
	Owned       bool         `json:"owned"`
	OwnerID     uint64       `json:"owner_id"`
	Tags        []string     `json:"tags"`
	TargetID    string       `json:"target_id"`
	Type        string       `json:"type"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Visibility  Visibility   `json:"visibility"`
}

// Validate returns an error if a constraint on the Object is not full-filled.
func (o *Object) Validate() error {
	if len(o.Attachments) > 5 {
		return ErrInvalidObject
	}

	for _, a := range o.Attachments {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	if o.OwnerID == 0 {
		return ErrInvalidObject
	}

	if len(o.Tags) > 5 {
		return ErrInvalidObject
	}

	if o.Type == "" {
		return ErrInvalidObject
	}

	if o.Visibility < 10 || o.Visibility > 50 {
		return ErrInvalidObject
	}

	return nil
}

// QueryOptions are passed to narrow down query for objects.
type QueryOptions struct {
	Deleted      bool
	ID           *uint64
	ObjectIDs    []uint64
	OwnerIDs     []uint64
	Owned        *bool
	Types        []string
	Visibilities []Visibility
}

// Service for object interactions.
type Service interface {
	Put(namespace string, object *Object) (*Object, error)
	Query(namespace string, opts QueryOptions) ([]*Object, error)
	Remove(namespace string, id uint64) error
	Setup(namespace string) error
	Teardown(namespace string) error
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// Visibility determines the visibility of Objects when consumed.
type Visibility uint8

func flakeNamespace(ns string) string {
	return fmt.Sprintf("%s_%s", ns, "objects")
}
