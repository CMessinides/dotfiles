package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type Renderer interface {
	RenderDocsetList(list DocsetList) error
	RenderEntryManifest(manifest EntryManifest) error
	RenderDocument(doc Document) error
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

func (r *ConsoleRenderer) RenderDocsetList(list DocsetList) error {
	w, err := r.text("Docsets")
	if err != nil {
		return err
	}
	defer w.Close()

	for _, d := range list.Docsets {
		_, err := fmt.Fprintf(w, "%s (%s)\n", d.Slug, d.FullName())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ConsoleRenderer) RenderEntryManifest(manifest EntryManifest) error {
	w, err := r.text("Entries: " + manifest.Docset)
	if err != nil {
		return err
	}
	defer w.Close()

	for _, e := range manifest.Entries {
		_, err := fmt.Fprintln(w, e.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ConsoleRenderer) RenderDocument(doc Document) error {
	filename := doc.Path
	if doc.Fragment != "" {
		filename += "#" + doc.Fragment
	}

	w, err := r.out(
		PagerVars{
			Filename: filename,
			Language: doc.Type.Lang(),
		},
	)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, doc.Reader())
	return err
}

func (r *ConsoleRenderer) text(header string) (io.WriteCloser, error) {
	return r.out(PagerVars{
		Filename: header,
		Language: "txt",
	})
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

func (r *PorcelainRenderer) RenderDocsetList(list DocsetList) error {
	for _, d := range list.Docsets {
		_, err := fmt.Fprintf(r.w, "%s\t%s\t%s\t%s\n", d.FullName(), d.Slug, d.Name, d.Release)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PorcelainRenderer) RenderEntryManifest(manifest EntryManifest) error {
	for _, e := range manifest.Entries {
		_, err := fmt.Fprintf(r.w, "%s\t%s\t%s\n", e.Name, e.Path, e.Type)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PorcelainRenderer) RenderDocument(doc Document) error {
	_, err := io.Copy(r.w, doc.Reader())
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

func (r *JSONRenderer) RenderDocsetList(list DocsetList) error {
	return r.e.Encode(list)
}

func (r *JSONRenderer) RenderEntryManifest(manifest EntryManifest) error {
	return r.e.Encode(manifest)
}

func (r *JSONRenderer) RenderDocument(doc Document) error {
	return r.e.Encode(doc)
}
