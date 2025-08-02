// Package report is a set of tools for richer error reporting in Go applications.
package report

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"strings"
)

// A Header encodes contextual information about an [Err]. The format is
// unspecified, but things like error codes, source locations, and similar
// metadata are well-suited for the header.
type Header interface {
	Header() string
}

// A Body encodes the main message of an [Err]. The format is unspecified, but
// the body should contain a human-readable message explaining the error.
type Body interface {
	Body() string
}

// Err is an error with a two-part message: a [Header] and a [Body].
type Err interface {
	error
	Header
	Body
}

type msgBuilder struct {
	msg    []byte
	prefix []byte // A buffer of data to prepend to msg.
	suffix []byte // A buffer of data to append to msg.
}

// Prepend adds the given bytes to the start of the message.
func (m *msgBuilder) Prepend(prefix []byte) {
	if m.prefix == nil {
		m.prefix = prefix
	} else {
		m.prefix = append(prefix, m.prefix...)
	}
}

// PrependString adds the given string to the start of the message.
func (m *msgBuilder) PrependString(prefix string) {
	m.Prepend([]byte(prefix))
}

// Append adds the given bytes to the end of the message.
func (m *msgBuilder) Append(suffix []byte) {
	if m.suffix == nil {
		m.suffix = suffix
	} else {
		m.suffix = append(m.suffix, suffix...)
	}
}

// AppendString adds the given string to the end of the message.
func (m *msgBuilder) AppendString(suffix string) {
	m.Append([]byte(suffix))
}

// Replace overwrites the current message with the given bytes.
func (m *msgBuilder) Replace(msg []byte) {
	m.msg = msg
	m.prefix = nil
	m.suffix = nil
}

// ReplaceString overwrites the current message with the given string.
func (m *msgBuilder) ReplaceString(msg string) {
	m.Replace([]byte(msg))
}

// String returns the complete message held by the builder.
func (m *msgBuilder) String() string {
	// Combine the buffers if needed.
	if len(m.prefix) > 0 || len(m.suffix) > 0 {
		buf := new(bytes.Buffer)
		buf.Grow(len(m.prefix) + len(m.msg) + len(m.suffix))
		buf.Write(m.prefix)
		buf.Write(m.msg)
		buf.Write(m.suffix)
		m.Replace(buf.Bytes())
	}

	return string(m.msg)
}

// An ErrorBuilder is used by [Middleware] to construct the components
// of an [Err].
type ErrorBuilder struct {
	Header *msgBuilder
	Body   *msgBuilder
}

func NewErrorBuilder(body string) *ErrorBuilder {
	h := new(msgBuilder)
	b := &msgBuilder{
		msg: []byte(body),
	}

	return &ErrorBuilder{
		Header: h,
		Body:   b,
	}
}

// Middleware applies a transformation to the header and body of an
// [Err].
type Middleware interface {
	Chain(...Middleware) Middleware
	Apply(*ErrorBuilder)
}

// MiddlewareFunc is a function that implements the [Middleware]
// interface.
type MiddlewareFunc func(*ErrorBuilder)

// Chain implements the [Middleware] interface.
func (mf MiddlewareFunc) Chain(m ...Middleware) Middleware {
	if len(m) == 0 {
		return mf
	}

	return MiddlewareFunc(func(b *ErrorBuilder) {
		mf.Apply(b)
		for _, next := range m {
			next.Apply(b)
		}
	})
}

// Apply implements the [Middleware] interface by calling itself on `err`.
func (mf MiddlewareFunc) Apply(b *ErrorBuilder) {
	mf(b)
}

// WithPrefix returns an [Middleware] that prepends the given prefix to the
// header of an [Err].
func WithPrefix(prefix string) Middleware {
	return MiddlewareFunc(func(b *ErrorBuilder) {
		b.Header.PrependString(prefix)
	})
}

// WithSuffix returns an [Middleware] that appends the given suffix to the
// header of an [Err].
func WithSuffix(suffix string) Middleware {
	return MiddlewareFunc(func(b *ErrorBuilder) {
		b.Header.AppendString(suffix)
	})
}

// WithFunctionLabel returns an [Middleware] that formats the given name and
// args as a function call and appends it to the header of an [Err].
//
// The formatting of each arg depends on its type:
//   - Bools, bytes, runes, and any number types are formatted using the
//     default formatter (equivalent to the `"%v"` verb from the `fmt`
//     package).
//   - Strings are formatted as quoted strings (`"%q"`). If the string is longer
//     than 20 characters, it will be truncated to 17 and `"..."` appended
//     before quoting.
//   - Every other type is formatted as its type name wrapped in angle brackets
//     (`"<%T>"`).
//
// For example, the transformer returned by the following call:
//
//	WithFunctionLabel("myFunc", 1, "example", []byte("example"))
//
// Is equivalent to:
//
//	WithSuffix(`myFunc(1, "example", <[]byte>)`)
func WithFunctionLabel(name string, args ...any) Middleware {
	s := new(strings.Builder)
	fmt.Fprintf(s, "%s(", name)
	for i, a := range args {
		switch v := a.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32,
			uint64, uintptr, float32, float64, complex64, complex128, bool:
			fmt.Fprint(s, v)
		case string:
			if len(v) > 20 {
				fmt.Fprintf(s, "%q", v[:17]+"...")
			} else {
				fmt.Fprintf(s, "%q", v)
			}
		default:
			fmt.Fprintf(s, "<%T>", v)
		}

		isLast := i == len(args)-1
		if isLast {
			fmt.Fprint(s, ")")
		} else {
			fmt.Fprint(s, ", ")
		}
	}

	return WithSuffix(s.String())
}

// WithMethodLabel returns an [Middleware] that formats the given name and args
// as a function call, prepends a literal ".", and appends the whole string to
// the header of an [Err]. See [WithFunctionLabel] for formatting details.
func WithMethodLabel(method string, args ...any) Middleware {
	return WithSuffix(".").Chain(WithFunctionLabel(method, args...))
}

// reportError is the internal implementation of [Err].
type reportError struct {
	header string
	body   string
	inner  error
}

// Header implements the [Err] interface.
func (r *reportError) Header() string {
	return r.header
}

// Body implements the [Err] interface.
func (r *reportError) Body() string {
	return r.body
}

// Error implements the [Err] interface.
func (r *reportError) Error() string {
	if r.inner != nil {
		return fmt.Sprintf("%s: %s: %s", r.header, r.body, r.inner.Error())
	} else {
		return fmt.Sprintf("%s: %s", r.header, r.body)
	}
}

// Unwrap returns the inner error, if any, for `errors.Unwrap()`.
func (r *reportError) Unwrap() error {
	return r.inner
}

// Newer is the interface that wraps the `New()` method for creating an [Err].
type Newer interface {
	// New constructs an [Err] with the given body.
	New(body string) Err
}

// Wrapper is the interface that wraps the `Wrap()` method for creating an
// [Err] that wraps another error.
type Wrapper interface {
	// Wrap constructs an [Err] with the given body and wraps it around the
	// given error.
	Wrap(body string, err error) Err
}

// A Chain is a sequence of middleware that can be applied repeatedly to create
// new errors that have similar structures.
type Chain struct {
	head Middleware
}

// NewChain returns a [Chain] that applies the given middleware.
func NewChain(middleware ...Middleware) *Chain {
	chain := new(Chain)

	if len(middleware) > 0 {
		head, rest := middleware[0], middleware[1:]
		chain.head = head.Chain(rest...)
	}

	return chain
}

// New implements the [Newer] interface. After the [Err] is created, the chain
// of middleware is applied to it.
func (c *Chain) New(body string) Err {
	return c.Wrap(body, nil)
}

// Wrap implements the [Wrapper] interface. After the [Err] is created, the
// chain of middleware is applied to it.
func (c *Chain) Wrap(body string, err error) Err {
	b := NewErrorBuilder(body)

	if c.head != nil {
		c.head.Apply(b)
	}

	return &reportError{
		header: b.Header.String(),
		body:   b.Body.String(),
		inner:  err,
	}
}

// Extend returns a new [Chain] that applies all the middleware of the original
// chain, followed by all the new middleware.
//
// If no middleware are provided (`len(m) == 0`), Extend returns the original
// chain.
func (c *Chain) Extend(m ...Middleware) *Chain {
	if len(m) == 0 {
		return c
	}

	return &Chain{
		head: c.head.Chain(m...),
	}
}

// Stack returns an iterator over header-body pairs in an error, starting
// with the outermost error and progressing through the chain of errors
// obtained by recursively calling [errors.Unwrap].
//
// For each error, if it implements the [Err] interface, the header is the
// return value of the `Header()` method, and the body the return value of
// the `Body()` method. Otherwise, the header is the type of the error
// wrapped in angle brackets (equivalent to the format string `"<%T>"`), and
// the body is the return value of the `Error()` method.
func Stack(err error) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		var r Err
		for err != nil {
			var header, body string
			if errors.As(err, &r) {
				header = r.Header()
				body = r.Body()
			} else {
				header = fmt.Sprintf("<%T>", err)
				body = err.Error()
			}

			if !yield(header, body) {
				return
			}

			err = errors.Unwrap(err)
		}
	}
}
