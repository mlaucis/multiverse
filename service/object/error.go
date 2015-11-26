package object

import "errors"

// Common errors for Object validation and Service.
var (
	ErrInvalidAttachment = errors.New("invalid attachment")
	ErrInvalidObject     = errors.New("invalid object")
	ErrNamespaceNotFound = errors.New("namespace not found")
)
