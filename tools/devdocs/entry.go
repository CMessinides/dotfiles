package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
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
}

func NewEntryIndex(entries []Entry) *EntryIndex {
	idx := &EntryIndex{
		paths:   make([]string, len(entries)),
		entries: make(map[string]*Entry),
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

// MarshalText implements the MarshalText method of the [Index] interface.
func (e *EntryIndex) MarshalText() (text []byte, err error) {
	buf := new(bytes.Buffer)

	for _, entry := range e.entries {
		fmt.Fprintf(buf, "%s\t%s\t%s\n", entry.Path, entry.Type, entry.Name)
	}

	return buf.Bytes(), nil
}

// UnmarshalText implements the UnmarshalText method of the [Index] interface.
func (e *EntryIndex) UnmarshalText(text []byte) error {
	entries := make(map[string]*Entry)
	scanner := bufio.NewScanner(bytes.NewReader(text))

	var n int
	for scanner.Scan() {
		n++
		line := scanner.Text()

		// Skip blank lines.
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			return NewErrBadIndexFormat(n, "not enough values (expected format <path>\\t<type>\\t<name>)")
		}

		path := parts[0]
		if path == "" {
			return NewErrBadIndexFormat(n, "path cannot be blank")
		}

		name := parts[2]
		if name == "" {
			return NewErrBadIndexFormat(n, "name cannot be blank")
		}

		entries[path] = &Entry{
			Path: path,
			Type: parts[1],
			Name: name,
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	e.entries = entries

	return nil
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
