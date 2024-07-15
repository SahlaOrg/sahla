package utils

import "fmt"

type ErrNotFound struct {
	Entity string
	ID     interface{}
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with ID %v not found", e.Entity, e.ID)
}

type ErrDuplicateEntry struct {
	Entity string
	Field  string
	Value  interface{}
}

func (e *ErrDuplicateEntry) Error() string {
	return fmt.Sprintf("duplicate entry for %s with %s: %v", e.Entity, e.Field, e.Value)
}

type ErrDatabase struct {
	Err error
}

func (e *ErrDatabase) Error() string {
	return fmt.Sprintf("database error: %v", e.Err)
}