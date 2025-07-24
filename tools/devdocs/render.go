package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Renderer interface {
	RenderDocsetList(docsets []Docset) error
	RenderEntryList(entries []*Entry) error
	RenderEntryView(view *EntryView) error
}

type ConsoleRenderer struct {
	stdout io.WriteCloser
	stderr io.WriteCloser
	isTTY  bool
}

func NewConsoleRenderer(stdout io.WriteCloser, stderr io.WriteCloser, isTTY bool) *ConsoleRenderer {
	return &ConsoleRenderer{
		stdout: stdout,
		stderr: stderr,
		isTTY:  isTTY,
	}
}

func (r *ConsoleRenderer) RenderDocsetList(docsets []Docset) error {
	w, err := r.text()
	if err != nil {
		return err
	}
	defer w.Close()

	for _, d := range docsets {
		_, err := fmt.Fprintf(w, "%s (%s)\n", d.Slug, d.FullName())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ConsoleRenderer) RenderEntryList(entries []*Entry) error {
	w, err := r.text()
	if err != nil {
		return err
	}
	defer w.Close()

	for _, e := range entries {
		_, err := fmt.Fprintln(w, e.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ConsoleRenderer) RenderEntryView(view *EntryView) error {
	filename := view.Document.Entry.String()

	w, err := r.out(
		PagerVars{
			Filename: filename,
			Language: "markdown",
		},
	)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = view.WriteTo(w)
	return err
}

func (r *ConsoleRenderer) text() (io.WriteCloser, error) {
	return r.out(PagerVars{})
}

func (r *ConsoleRenderer) out(vars PagerVars) (io.WriteCloser, error) {
	if !r.isTTY {
		return r.stdout, nil
	}

	p, err := LookupPager(vars)
	if err != nil {
		return nil, err
	}

	pw, err := p.Wrap(r.stdout, r.stderr)
	if err != nil {
		return nil, err
	}

	return pw, nil
}

type PorcelainRenderer struct {
	w io.Writer
}

func NewPorcelainRenderer(w io.Writer) *PorcelainRenderer {
	return &PorcelainRenderer{w: w}
}

func (r *PorcelainRenderer) RenderDocsetList(docsets []Docset) error {
	for _, d := range docsets {
		_, err := fmt.Fprintf(r.w, "%s\t%s\t%s\t%s\n", d.FullName(), d.Slug, d.Name, d.Release)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PorcelainRenderer) RenderEntryList(entries []*Entry) error {
	for _, e := range entries {
		_, err := fmt.Fprintf(r.w, "%s\t%s\t%s\n", e.Path, e.Type, e.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PorcelainRenderer) RenderEntryView(entry *EntryView) error {
	_, err := entry.WriteTo(r.w)
	return err
}

type JSONRenderer struct {
	e *json.Encoder
}

func NewJSONRenderer(w io.Writer) *JSONRenderer {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)

	return &JSONRenderer{
		e: e,
	}
}

func (r *JSONRenderer) RenderDocsetList(docsets []Docset) error {
	return r.e.Encode(docsets)
}

func (r *JSONRenderer) RenderEntryList(entries []*Entry) error {
	return r.e.Encode(entries)
}

func (r *JSONRenderer) RenderEntryView(view *EntryView) error {
	s := new(strings.Builder)
	_, err := view.WriteTo(s)
	if err != nil {
		return err
	}

	return r.e.Encode(struct {
		Docset  string       `json:"docset"`
		Entry   EntryLocator `json:"entry"`
		Lines   *LineRange   `json:"lines"`
		Content string       `json:"content"`
	}{
		Docset:  view.Document.Docset,
		Entry:   view.Document.Entry,
		Lines:   view.Lines,
		Content: s.String(),
	})
}
