package main

import (
	"fmt"
	"slices"
	"strings"
)

type StructuredError struct {
	Label   string
	Message string
	inner   error
}

// String implements the [fmt.Stringer] interface.
func (s *StructuredError) String() string {
	return s.Error()
}

// Error implements the error interface.
func (s *StructuredError) Error() string {
	m := fmt.Sprintf("%s: %s", s.Label, s.Message)
	if s.inner != nil {
		return m + ": " + s.inner.Error()
	} else {
		return m
	}
}

func (s *StructuredError) Unwrap() error {
	return s.inner
}

type ErrorTransformer interface {
	TransformError(err *StructuredError)
}

type ErrorTransformerFunc func(err *StructuredError)

func (e ErrorTransformerFunc) TransformError(err *StructuredError) {
	e(err)
}

type ErrorBuilder struct {
	transforms []ErrorTransformer
}

func (e *ErrorBuilder) Extend(t ...ErrorTransformer) *ErrorBuilder {
	return &ErrorBuilder{
		transforms: slices.Concat(e.transforms, t),
	}
}

func (e *ErrorBuilder) New(msg string) *StructuredError {
	return e.build(&StructuredError{
		Message: msg,
	})
}

func (e *ErrorBuilder) Wrap(msg string, err error) *StructuredError {
	return e.build(&StructuredError{
		Message: msg,
		inner:   err,
	})
}

func (e *ErrorBuilder) build(err *StructuredError) *StructuredError {
	for _, t := range e.transforms {
		t.TransformError(err)
	}

	return err
}

func NewErrorBuilder(t ...ErrorTransformer) *ErrorBuilder {
	return &ErrorBuilder{
		transforms: t,
	}
}

func WithPrefix(prefix string) ErrorTransformer {
	return ErrorTransformerFunc(func(err *StructuredError) {
		err.Label = prefix + err.Label
	})
}

func WithSuffix(suffix string) ErrorTransformer {
	return ErrorTransformerFunc(func(err *StructuredError) {
		err.Label = err.Label + suffix
	})
}

func WithMethodLabel(method string, args ...any) ErrorTransformer {
	s := new(strings.Builder)
	fmt.Fprintf(s, ".%s(", method)
	for i, a := range args {
		isLast := i == len(args)-1
		if isLast {
			fmt.Fprintf(s, "%q)", a)
		} else {
			fmt.Fprintf(s, "%q, ", a)
		}
	}

	return WithSuffix(s.String())
}
