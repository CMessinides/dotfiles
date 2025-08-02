package main

import (
	"encoding"
	"fmt"

	"github.com/cmessinides/dotfiles/tools/devdocs/internal/report"
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
	report.Err
	Line int
}

func NewBadIndexFormatError(r report.Newer, line int, msg string) *ErrBadIndexFormat {
	return &ErrBadIndexFormat{
		Err: r.New(
			fmt.Sprintf("failed to parse index line %d: %s", line, msg),
		),
		Line: line,
	}
}
