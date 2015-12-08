package object

import "errors"

// Common errors for Object validation and Service.
var (
	ErrInvalidAttachment = errors.New("invalid attachment")
	ErrInvalidObject     = errors.New("invalid object")
	ErrMissingReference  = errors.New("referenced object missing")
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrNotFound          = errors.New("object not found")
)
