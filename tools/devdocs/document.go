package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

type DocumentContent []byte

func (d DocumentContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

func (d DocumentContent) Reader() io.Reader {
	return bytes.NewReader(d)
}

type document struct {
	Content DocumentContent `json:"content"`
	Docset  string          `json:"docset"`
	Entry   EntryLocator    `json:"entry"`
}

type LineRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type DocumentSection struct {
	Level int       `json:"level"`
	ID    string    `json:"id"`
	Lines LineRange `json:"lines"`
}

type DocumentIndex struct {
	IDs    []string
	Ranges map[string]LineRange
}

func NewDocumentIndex(sections []*DocumentSection) *DocumentIndex {
	idx := &DocumentIndex{
		IDs:    make([]string, 0, len(sections)),
		Ranges: make(map[string]LineRange),
	}

	for _, s := range sections {
		if s.ID == "" {
			continue
		}

		idx.IDs = append(idx.IDs, s.ID)
		idx.Ranges[s.ID] = s.Lines
	}

	return idx
}

func BuildDocumentIndex(md []byte, ids []string) (*DocumentIndex, error) {
	errs := NewErrorBuilder(
		WithFunctionLabel("BuildDocumentIndex", md, ids),
	)
	sections := make([]*DocumentSection, 0, len(ids))
	scanner := bufio.NewScanner(bytes.NewReader(md))

	i := 0
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Text()

		// Skip non-heading lines.
		if !strings.HasPrefix(line, "#") {
			continue
		}

		// Count leading '#' characters.
		var lvl int
		for line[lvl] == '#' {
			lvl++
		}

		// Skip H1s.
		if lvl == 1 {
			continue
		}

		s := &DocumentSection{
			Level: lvl,
			ID:    ids[i],
			Lines: LineRange{
				Start: lineno,
				End:   -1,
			},
		}

		slog.Debug("found section", "level", s.Level, "id", s.ID, "line", s.Lines.Start)

		sections = append(sections, s)
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, errs.Wrap("encounted error while scanning Markdown", err)
	}

	slog.Debug("calculating section ranges", "count", len(sections))
	stack := newStack[*DocumentSection]()
	for _, cur := range sections {
		slog.Debug("processing section", "id", cur.ID, "stacked", stack.Len())
	inner:
		for stack.Len() > 0 {
			prev := stack.Top()
			if prev.Level >= cur.Level {
				stack.Pop()
				prev.Lines.End = cur.Lines.Start - 1
				slog.Debug("calculated range for section",
					"id", prev.ID,
					"start", prev.Lines.Start,
					"end", prev.Lines.End,
				)
			} else {
				break inner
			}
		}

		stack.Push(cur)
	}

	slog.Debug("cleaning up section stack", "count", stack.Len())
	for stack.Len() > 0 {
		rem, _ := stack.Pop()
		rem.Lines.End = lineno
		slog.Debug("calculated range for section",
			"id", rem.ID,
			"start", rem.Lines.Start,
			"end", rem.Lines.End,
		)
	}

	idx := NewDocumentIndex(sections)
	slog.Debug("created document index", "count", idx.Count())
	return idx, nil
}

// Get implements the Get method of the [Index] interface.
func (d *DocumentIndex) Get(id string) (lines *LineRange, ok bool) {
	r, ok := d.Ranges[id]
	if ok {
		lines = &r
	}

	return lines, ok
}

// MarshalText implements the MarshalText method of the [Index] interface.
func (d *DocumentIndex) MarshalText() (text []byte, err error) {
	buf := new(bytes.Buffer)

	for _, id := range d.IDs {
		l := d.Ranges[id]
		fmt.Fprintf(buf, "%d:%d %s\n", l.Start, l.End, id)
	}

	return buf.Bytes(), nil
}

// UnmarshalText implements the UnmarshalText method of the [Index] interface.
func (d *DocumentIndex) UnmarshalText(text []byte) error {
	ids := make([]string, 0)
	ranges := make(map[string]LineRange)
	scanner := bufio.NewScanner(bytes.NewReader(text))

	var n int
	for scanner.Scan() {
		n++
		line := scanner.Text()

		// Skip blank lines.
		if line == "" {
			continue
		}

		rawRange, id, ok := strings.Cut(line, " ")
		if !ok {
			return NewErrBadIndexFormat(n, "no space after range (expected format <start>:<end> <id>)")
		}

		rawStart, rawEnd, ok := strings.Cut(rawRange, ":")
		if !ok {
			return NewErrBadIndexFormat(n, "bad range syntax (expected format <start>:<end>)")
		}

		start, err := strconv.Atoi(rawStart)
		if err != nil {
			return NewErrBadIndexFormat(n, "bad range syntax: <start> must be a number")
		}

		end, err := strconv.Atoi(rawEnd)
		if err != nil {
			return NewErrBadIndexFormat(n, "bad range syntax : <end> must be a number")
		}

		ids = append(ids, id)
		ranges[id] = LineRange{
			Start: start,
			End:   end,
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	d.IDs = ids
	d.Ranges = ranges

	return nil
}

func (d *DocumentIndex) Count() int {
	return len(d.IDs)
}

// HTMLDocument is the HTML documentation for an entry, as retrieved from the
// DevDocs API.
type HTMLDocument struct {
	document
}

func NewHTMLDocument(docset string, entry EntryLocator, content []byte) *HTMLDocument {
	return &HTMLDocument{
		document: document{
			Docset:  docset,
			Entry:   entry,
			Content: DocumentContent(content),
		},
	}
}

type MarkdownDocument struct {
	document
	Index Index[*LineRange]
}

func NewMarkdownDocument(docset string, entry EntryLocator, content []byte, idx *DocumentIndex) *MarkdownDocument {
	return &MarkdownDocument{
		document: document{
			Docset:  docset,
			Entry:   entry,
			Content: DocumentContent(content),
		},
		Index: idx,
	}
}

func NewMarkdownDocumentFromHTML(html *HTMLDocument, md []byte, idx *DocumentIndex) *MarkdownDocument {
	return NewMarkdownDocument(html.Docset, html.Entry, md, idx)
}
