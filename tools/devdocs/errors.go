package main

import (
	"fmt"
	"slices"
	"strings"
)

type LabeledError interface {
	error
	Label() string
	Body() string
}

type StructuredError struct {
	label string
	msg   string
	inner error
}

// Label implements the [LabeledError] interface.
func (s *StructuredError) Label() string {
	return s.label
}

// Body implements the [LabeledError] interface.
func (s *StructuredError) Body() string {
	return s.msg
}

func (s *StructuredError) Unwrap() error {
	return s.inner
}

// Error implements the error interface.
func (s *StructuredError) Error() string {
	m := fmt.Sprintf("%s: %s", s.label, s.msg)
	if s.inner != nil {
		return m + ": " + s.inner.Error()
	} else {
		return m
	}
}

type ErrorTransformer interface {
	TransformError(err *StructuredError)
	Chain(t ...ErrorTransformer) ErrorTransformer
}

type ErrorTransformerFunc func(err *StructuredError)

func (e ErrorTransformerFunc) TransformError(err *StructuredError) {
	e(err)
}

func (e ErrorTransformerFunc) Chain(t ...ErrorTransformer) ErrorTransformer {
	return ErrorTransformerFunc(func(err *StructuredError) {
		e.TransformError(err)
		for _, next := range t {
			next.TransformError(err)
		}
	})
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
		msg: msg,
	})
}

func (e *ErrorBuilder) Wrap(msg string, err error) *StructuredError {
	return e.build(&StructuredError{
		msg:   msg,
		inner: err,
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
		err.label = prefix + err.label
	})
}

func WithSuffix(suffix string) ErrorTransformer {
	return ErrorTransformerFunc(func(err *StructuredError) {
		err.label = err.label + suffix
	})
}

func WithFunctionLabel(name string, args ...any) ErrorTransformer {
	s := new(strings.Builder)
	fmt.Fprintf(s, "%s(", name)
	for i, a := range args {
		isLast := i == len(args)-1
		if isLast {
			fmt.Fprintf(s, "%#v)", a)
		} else {
			fmt.Fprintf(s, "%#v, ", a)
		}
	}

	return WithSuffix(s.String())
}

func WithMethodLabel(method string, args ...any) ErrorTransformer {
	return WithSuffix(".").Chain(WithFunctionLabel(method, args...))
}
