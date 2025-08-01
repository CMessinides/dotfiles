package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/cmessinides/dotfiles/tools/devdocs/internal/report"
	"github.com/fatih/color"
)

type Reporter interface {
	ReportStatus(msg string)
	ReportWarning(msg string)
	ReportError(err error)
}

type colorMap struct {
	bold       *color.Color
	dim        *color.Color
	error      *color.Color
	warning    *color.Color
	arg        *color.Color
	cmd        *color.Color
	stackLabel *color.Color
}

var colors = colorMap{
	bold:       color.New(color.Bold),
	dim:        color.New(color.Faint),
	error:      color.New(color.FgRed, color.Bold),
	warning:    color.New(color.FgYellow, color.Bold),
	arg:        color.New(color.FgGreen),
	cmd:        color.New(color.FgCyan, color.Bold),
	stackLabel: color.New(color.FgBlue),
}

func quote(a string) string {
	return fmt.Sprintf(`"%s"`, a)
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

	for h, b := range report.Stack(err) {
		fmt.Fprintf(s, "%s: %s\n", colors.stackLabel.Sprint(h), b)
	}

	return s.String()
}

//go:embed templates/*.tmpl
var tmplFS embed.FS

var funcMap = template.FuncMap{
	"cmd":     colors.cmd.Sprint,
	"arg":     colors.arg.Sprint,
	"error":   colors.error.Sprint,
	"warning": colors.warning.Sprint,
	"bold":    colors.bold.Sprint,
	"dim":     colors.dim.Sprint,
	"quote":   quote,
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
		Code: "default",
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
