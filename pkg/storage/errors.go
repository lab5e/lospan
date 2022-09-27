package storage

import "errors"

// ErrNotFound is returned by the storage layer when the specified entity
// isn't found.
var ErrNotFound = errors.New("not found")

// ErrDeleteConstraint is returned by the storage layer when there's an
// constraint preventing the entity from being deleted
var ErrDeleteConstraint = errors.New("constraint error")

// ErrAlreadyExists is returned by the storage layer when the item already
// exists, ie there's a duplicate
var ErrAlreadyExists = errors.New("already exists")
