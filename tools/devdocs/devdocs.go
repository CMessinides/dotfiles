package main

import (
	"bytes"
	"encoding/json"
	"io"
)

// Docset represents a docset retrieved from the DevDocs API.
type Docset struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Release string `json:"release"`
}

func (d Docset) FullName() string {
	if d.Release == "" {
		return d.Name
	}

	return d.Name + " " + d.Release
}

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

type ContentType int

const (
	ContentText ContentType = iota
	ContentHTML
	ContentMarkdown
)

func (c ContentType) String() string {
	switch c {
	case ContentHTML:
		return "html"
	case ContentMarkdown:
		return "markdown"
	default:
		return "text"
	}
}

func (c ContentType) Lang() string {
	switch c {
	case ContentHTML:
		return "html"
	case ContentMarkdown:
		return "markdown"
	default:
		return "txt"
	}
}

func (c ContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// Document is the text documentation for an entry.
type Document struct {
	Content  []byte
	Type     ContentType
	Docset   string
	Path     string
	Fragment string
}

var EmptyDocument = Document{}

func NewDocument(docset string, path string, fragment string, contentType ContentType, content []byte) Document {
	return Document{
		Content:  content,
		Type:     contentType,
		Docset:   docset,
		Path:     path,
		Fragment: fragment,
	}
}

func (d Document) Reader() io.Reader {
	return bytes.NewReader(d.Content)
}

func (d Document) ToMarkdown(md []byte) Document {
	return Document{
		Content:  md,
		Type:     ContentMarkdown,
		Docset:   d.Docset,
		Path:     d.Path,
		Fragment: d.Fragment,
	}
}

func (d Document) MarshalJSON() ([]byte, error) {
	c := struct {
		Docset      string      `json:"docset"`
		Path        string      `json:"path"`
		Fragment    string      `json:"fragment"`
		ContentType ContentType `json:"contentType"`
		Content     string      `json:"content"`
	}{
		Docset:      d.Docset,
		Path:        d.Path,
		Fragment:    d.Fragment,
		ContentType: d.Type,
		Content:     string(d.Content),
	}

	return json.Marshal(c)
}
