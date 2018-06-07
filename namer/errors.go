package namer

import "errors"

var (
	// ErrNotFound is returned by Get() or Update() when the resource was not found by ID.
	ErrNotFound = errors.New("resource was not found by ID or name")
)
