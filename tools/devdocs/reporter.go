package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type Reporter interface {
	ReportStatus(msg string)
	ReportWarning(msg string)
	ReportError(err error)
}

func fmtBold(a ...any) string {
	return color.New(color.Bold).Sprint(a...)
}

func fmtDim(a ...any) string {
	return color.New(color.Faint).Sprint(a...)
}

func fmtError(a ...any) string {
	return color.New(color.FgRed, color.Bold).Sprint(a...)
}

func fmtWarning(a ...any) string {
	return color.New(color.FgYellow, color.Bold).Sprint(a...)
}

func fmtQuote(a string) string {
	return fmt.Sprintf(`"%s"`, a)
}

func fmtArg(a ...any) string {
	s := color.New(color.FgYellow).Sprint(a...)
	if color.NoColor {
		s = fmtQuote(s)
	}

	return s
}

func fmtDocset(a ...any) string {
	return color.New(color.FgMagenta).Sprint(a...)
}

func fmtEntry(a ...any) string {
	return color.New(color.FgBlue).Sprint(a...)
}

func fmtCmd(a ...any) string {
	return color.New(color.FgCyan, color.Bold).Sprint(a...)
}

func indent(count int, text string) string {
	s := new(strings.Builder)

	spaces := strings.Repeat(" ", count)
	for line := range strings.Lines(text) {
		fmt.Fprint(s, spaces+line)
	}

	return s.String()
}

func inspect(err error) string {
	s := new(strings.Builder)

	var l LabeledError
	for err != nil {
		if errors.As(err, &l) {
			fmt.Fprintf(s, "%s: %s\n", l.Label(), l.Body())
		} else {
			fmt.Fprintf(s, "<err>: %s\n", err.Error())
		}

		err = errors.Unwrap(err)
	}

	return s.String()
}

//go:embed templates/*.tmpl
var tmplFS embed.FS

var funcMap = template.FuncMap{
	"cmd":     fmtCmd,
	"arg":     fmtArg,
	"docset":  fmtDocset,
	"entry":   fmtEntry,
	"error":   fmtError,
	"warning": fmtWarning,
	"bold":    fmtBold,
	"dim":     fmtDim,
	"quote":   fmtQuote,
	"indent":  indent,
	"inspect": inspect,
}

var templates = template.Must(
	template.
		New("").
		Funcs(funcMap).
		ParseFS(mustSubFS(tmplFS, "templates"), "*.tmpl"),
)

type ConsoleReporter struct {
	w       io.Writer
	t       *template.Template
	Verbose bool
}

func NewConsoleReporter(w io.Writer) *ConsoleReporter {
	return &ConsoleReporter{
		w: w,
		t: templates,
	}
}

type ErrorReport struct {
	Code    string
	Err     error
	Full    error
	Verbose bool
}

func (r *ErrorReport) narrow(code string, err error) *ErrorReport {
	r.Code = code
	r.Err = err
	return r
}

func DefaultErrorReport(err error) *ErrorReport {
	return &ErrorReport{
		Code: "generic",
		Err:  err,
		Full: err,
	}
}

func NewErrorReport(err error) *ErrorReport {
	r := DefaultErrorReport(err)

	var dnf *DocsetNotFoundError
	if errors.As(err, &dnf) {
		return r.narrow("docset-not-found", dnf)
	}

	var enf *EntryNotFoundError
	if errors.As(err, &enf) {
		return r.narrow("entry-not-found", enf)
	}

	return r
}

func (c *ConsoleReporter) ReportError(err error) {
	r := NewErrorReport(err)
	r.Verbose = c.Verbose

	t := fmt.Sprintf("err-%s.tmpl", r.Code)
	c.executeTemplate(t, r)
}

func (c *ConsoleReporter) ReportWarning(msg string) {
}

func (c *ConsoleReporter) ReportStatus(msg string) {
}

func (c *ConsoleReporter) executeTemplate(name string, data any) error {
	err := c.t.ExecuteTemplate(c.w, name, data)
	if err != nil {
		fmt.Fprintf(c.w, "error rendering template %q: %s\n", name, err)
	}

	return err
}

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
