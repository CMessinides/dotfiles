package main

import (
	"encoding"
	"fmt"
)

// Index is a serializable data structure that facilitates fast lookups with
// string keys. Think of it as a readonly Go map that can be converted to and
// from plain text.
type Index[T any] interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	Get(key string) (value T, ok bool)
}

type ErrBadIndexFormat struct {
	Line int
	msg  string
}

func NewErrBadIndexFormat(line int, msg string) *ErrBadIndexFormat {
	return &ErrBadIndexFormat{
		Line: line,
		msg:  msg,
	}
}

func (e *ErrBadIndexFormat) Error() string {
	return fmt.Sprintf("failed to parse index line %d: %s", e.Line, e.msg)
}
