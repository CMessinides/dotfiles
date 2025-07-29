package main

import (
	"context"
	"fmt"
)

type EntryNotFoundError struct {
	Docset string
	Path   string
}

// Error implements the error interface.
func (e *EntryNotFoundError) Error() string {
	return fmt.Sprintf("entry %q not found in docset %q", e.Path, e.Docset)
}

type SectionNotFoundError struct {
	ID     string
	Path   string
	Docset string
}

// Error implements the error interface.
func (s *SectionNotFoundError) Error() string {
	return fmt.Sprintf("section %q not found in document %q", s.ID, s.Path)
}

type Service struct {
	cache     Cache
	client    *Client
	converter *MarkdownConverter
}

func NewService(cache Cache, client *Client, converter *MarkdownConverter) *Service {
	return &Service{
		cache:     cache,
		client:    client,
		converter: converter,
	}
}

func (s *Service) ListDocsets(ctx context.Context) ([]Docset, error) {
	fail := func(err error, msg string) error {
		return fmt.Errorf("ListDocsets(): %s: %w", msg, err)
	}

	d, err := s.client.ListDocsets(ctx)
	if err != nil {
		return d, fail(err, "could not get docsets")
	}

	return d, nil
}

func (s *Service) ListEntries(ctx context.Context, docset string) ([]*Entry, error) {
	fail := func(err error, msg string) error {
		return fmt.Errorf("service.ListEntries(%q): %s: %w", docset, msg, err)
	}

	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, fail(err, "could not index entries")
	}

	return idx.Entries(), nil
}

func (s *Service) ShowEntry(ctx context.Context, docset string, path string) (*EntryView, error) {
	fail := func(err error, msg string) error {
		return fmt.Errorf("service.ShowEntry(%q, %q): %s: %w", docset, path, msg, err)
	}

	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, fail(err, "could not index entries")
	}

	entry, ok := idx.Get(path)
	if !ok {
		return nil, fail(&EntryNotFoundError{
			Docset: docset,
			Path:   path,
		}, "entry not found")
	}

	loc := NewEntryLocator(entry.Path)
	html, err := s.client.GetDocument(ctx, docset, loc)
	if err != nil {
		return nil, fail(err, "could not get document")
	}

	md, err := s.converter.Convert(html)
	if err != nil {
		return nil, fail(err, "could not convert document to Markdown")
	}

	var view *EntryView
	if loc.HasFragment() {
		lines, ok := md.Index.Get(loc.Fragment)
		if !ok {
			return nil, fail(&SectionNotFoundError{
				ID:     loc.Fragment,
				Path:   loc.Path,
				Docset: docset,
			}, "section not found")
		}

		view = NewExcerptView(md, lines)
	} else {
		view = NewDocumentView(md)
	}

	return view, nil
}

func (s *Service) entryIndex(ctx context.Context, docset string) (*EntryIndex, error) {
	fail := func(err error, msg string) error {
		return fmt.Errorf("service.entryIndex(%q): %s: %w", docset, msg, err)
	}

	m, err := s.client.ListEntries(ctx, docset)
	if err != nil {
		return nil, fail(err, "could not get entries")
	}

	return NewEntryIndex(m.Entries), nil
}
