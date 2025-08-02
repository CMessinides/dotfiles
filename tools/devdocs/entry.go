package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/cmessinides/dotfiles/tools/devdocs/internal/report"
)

// Entry represents a documentation entry retrieved from the DevDocs API.
// Docsets contain multiple entries -- you can think of them as pages.
type Entry struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// EntryManifest is the collection of entries within a docset, as retrieved
// from the DevDocs API.
type EntryManifest struct {
	Docset  string  `json:"docset"`
	Entries []Entry `json:"entries"`
}

type EntryIndex struct {
	paths   []string
	entries map[string]*Entry
	errs    *report.Chain
}

func NewEntryIndex(entries []Entry) *EntryIndex {
	idx := &EntryIndex{
		paths:   make([]string, len(entries)),
		entries: make(map[string]*Entry),
		errs:    report.NewChain(report.WithPrefix("EntryIndex")),
	}

	for i, e := range entries {
		idx.paths[i] = e.Path
		idx.entries[e.Path] = &e
	}

	return idx
}

func (e *EntryIndex) Entries() []*Entry {
	entries := make([]*Entry, len(e.paths))
	for i, p := range e.paths {
		entries[i] = e.entries[p]
	}

	return entries
}

// Get implements the Get method of the [Index] interface.
func (e *EntryIndex) Get(path string) (entry *Entry, ok bool) {
	entry, ok = e.entries[path]
	if ok {
		return
	}

	// If the path didn't work as-is, try adding "/index" to the path.
	loc := NewEntryLocator(path)
	loc.Path += "/index"
	entry, ok = e.entries[loc.String()]

	return entry, ok
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (e *EntryIndex) MarshalText() (text []byte, err error) {
	// errs := e.errs.Extend(report.WithMethodLabel("MarshalText"))
	buf := new(bytes.Buffer)

	for _, entry := range e.entries {
		fmt.Fprintf(buf, "%s\t%s\t%s\n", entry.Path, entry.Type, entry.Name)
	}

	return buf.Bytes(), nil
}

// WriteTo implements the [io.WriterTo] interface.
func (e *EntryIndex) WriteTo(w io.Writer) (n int64, err error) {
	errs := e.errs.Extend(report.WithMethodLabel("WriteTo", w))

	var i int
	for _, entry := range e.entries {
		ln, err := fmt.Fprintf(w, "%s\t%s\t%s\n", entry.Path, entry.Type, entry.Name)
		if err != nil {
			return n, errs.Wrap(
				fmt.Sprintf("failed to write entry %q (line %d)", entry.Path, i),
				err,
			)
		}

		i++
		n += int64(ln)
	}

	return n, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (e *EntryIndex) UnmarshalText(text []byte) error {
	errs := e.errs.Extend(report.WithMethodLabel("UnmarshalText", text))
	_, err := e.ReadFrom(bytes.NewReader(text))
	if err != nil {
		return errs.Wrap("failed to read index", err)
	}

	return nil
}

// ReadFrom implements the [io.ReaderFrom] interface.
func (e *EntryIndex) ReadFrom(r io.Reader) (n int64, err error) {
	errs := e.errs.Extend(report.WithMethodLabel("ReadFrom", r))
	entries := make(map[string]*Entry)
	scanner := bufio.NewScanner(r)

	var l int
	for scanner.Scan() {
		l++
		n += int64(len(scanner.Bytes()) + 1)
		line := scanner.Text()

		// Skip blank lines.
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			return 0, NewBadIndexFormatError(errs, l, "not enough values (expected format <path>\\t<type>\\t<name>)")
		}

		path := parts[0]
		if path == "" {
			return 0, NewBadIndexFormatError(errs, l, "path cannot be blank")
		}

		name := parts[2]
		if name == "" {
			return 0, NewBadIndexFormatError(errs, l, "name cannot be blank")
		}

		entries[path] = &Entry{
			Path: path,
			Type: parts[1],
			Name: name,
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, errs.Wrap("failed to scan from reader", err)
	}

	e.entries = entries

	return n, nil
}

// EntryLocator encodes how to locate an entry from DevDocs, with the
// path to a document and, optionally, a fragment identifying a section
// within that document.
type EntryLocator struct {
	Path     string `json:"path"`
	Fragment string `json:"fragment"`
}

func NewEntryLocator(rawPath string) EntryLocator {
	path, frag, _ := strings.Cut(rawPath, "#")
	return EntryLocator{
		Path:     path,
		Fragment: frag,
	}
}

func (e EntryLocator) HasFragment() bool {
	return e.Fragment != ""
}

// String implements the [fmt.Stringer] interface.
func (e EntryLocator) String() string {
	s := e.Path
	if e.HasFragment() {
		s += "#" + e.Fragment
	}

	return s
}

type EntryView struct {
	Lines    *LineRange
	Document *MarkdownDocument
}

func NewExcerptView(doc *MarkdownDocument, lines *LineRange) *EntryView {
	return &EntryView{
		Lines:    lines,
		Document: doc,
	}
}

func NewDocumentView(doc *MarkdownDocument) *EntryView {
	return &EntryView{
		Document: doc,
	}
}

func (e *EntryView) IsExcerpt() bool {
	return e.Lines != nil
}

// WriteTo implements the [io.WriterTo] interface.
func (e *EntryView) WriteTo(w io.Writer) (n int64, err error) {
	if !e.IsExcerpt() {
		return io.Copy(w, e.Document.Content.Reader())
	} else {
		scanner := bufio.NewScanner(e.Document.Content.Reader())
		start := e.Lines.Start
		end := e.Lines.End

		var line, lineN int
		for scanner.Scan() {
			line++

			if line < start {
				continue
			} else if line > end {
				break
			}

			b := append(scanner.Bytes(), '\n')
			lineN, err = w.Write(b)
			n += int64(lineN)
			if err != nil {
				return n, err
			}
		}

		if err := scanner.Err(); err != nil {
			return n, err
		}

		return n, err
	}
}
