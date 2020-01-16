package main

import (
	"errors"
	"net/http"
)

var (
	// ErrNotFound is a 404 error.
	ErrNotFound = NewError(http.StatusNotFound, errors.New("not found"))
	// ErrMethodNotAllowed is a 405 error.
	ErrMethodNotAllowed = NewError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
)

// Error is a HTTP error.
type Error interface {
	error
	Status() int
}

type httpError struct {
	error
	status int
}

func (e *httpError) Status() int { return e.status }

// NewError creates a new HTTP error.
func NewError(status int, err error) Error {
	return &httpError{err, status}
}
