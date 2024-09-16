package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	ErrCantDoIt       = errors.New("cant do it")
)
