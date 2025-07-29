package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type Reporter interface {
	ReportStatus(msg string)
	ReportWarning(msg string)
	ReportError(err error)
}

func bold(a ...any) string {
	return color.New(color.Bold).Sprint(a...)
}

func dim(a ...any) string {
	return color.New(color.Faint).Sprint(a...)
}

func red(a ...any) string {
	return color.New(color.FgRed).Sprint(a...)
}

func yellow(a ...any) string {
	return color.New(color.FgYellow).Sprint(a...)
}

func cyan(a ...any) string {
	return color.New(color.FgCyan).Sprint(a...)
}

func quote(a any) string {
	return fmt.Sprintf(`"%s"`, a)
}

func emph(a any) string {
	if color.NoColor {
		return quote(a)
	} else {
		return yellow(a)
	}
}

func indent(a any) string {
	return fmt.Sprintf("  %s", a)
}

var funcPrefixPattern = regexp.MustCompile(`^\s*(?:\w+\.)?\w+\([^\)]*\)`)

func printstack(err error) string {
	s := new(strings.Builder)
	parts := strings.SplitSeq(err.Error(), ":")

	for p := range parts {
		if funcPrefixPattern.MatchString(p) {
			s.WriteString("\n  " + cyan(strings.TrimSpace(p)))
		} else {
			s.WriteString(dim(":" + p))
		}
	}

	return s.String()
}

//go:embed templates/*.tmpl
var tmplFS embed.FS

var funcMap = template.FuncMap{
	"bold":       bold,
	"dim":        dim,
	"red":        red,
	"yellow":     yellow,
	"cyan":       cyan,
	"quote":      quote,
	"emph":       emph,
	"indent":     indent,
	"printstack": printstack,
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

type viewData struct {
	Verbose bool
	Err     error
	Full    error
	subject string
}

func (e viewData) Subject() string {
	if e.subject != "" {
		return e.subject
	}

	return e.Err.Error()
}

func (c *ConsoleReporter) ReportError(err error) {
	var dnf *DocsetNotFoundError
	if errors.As(err, &dnf) {
		c.reportDocsetNotFound(dnf, err)
		return
	}

	var enf *EntryNotFoundError
	if errors.As(err, &enf) {
		c.reportEntryNotFound(enf, err)
		return
	}

	c.reportGenericError(err)
}

func (c *ConsoleReporter) ReportWarning(msg string) {
}

func (c *ConsoleReporter) ReportStatus(msg string) {
}

func (c *ConsoleReporter) reportDocsetNotFound(err *DocsetNotFoundError, full error) {
	c.error(
		err,
		full,
		fmt.Sprintf("docset %s not found", emph(err.Docset)),
		"err-docset-not-found.tmpl",
	)
}

func (c *ConsoleReporter) reportEntryNotFound(err *EntryNotFoundError, full error) {
	c.error(
		err,
		full,
		fmt.Sprintf(
			"entry %s not found in docset %s",
			emph(err.Path),
			emph(err.Docset),
		),
		"err-entry-not-found.tmpl",
	)
}

func (c *ConsoleReporter) reportGenericError(err error) {
	c.error(err, err, "an unknown error occurred", "err-generic.tmpl")
}

func (c *ConsoleReporter) error(err error, full error, subject string, template string) {
	d := viewData{
		Verbose: c.Verbose,
		Err:     err,
		Full:    full,
		subject: subject,
	}

	c.executeTemplate("err-header.tmpl", d)
	c.executeTemplate(template, d)
	c.executeTemplate("err-footer.tmpl", d)
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
