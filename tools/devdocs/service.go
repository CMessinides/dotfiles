package main

import (
	"context"
	"fmt"
)

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
	return s.client.ListDocsets(ctx)
}

func (s *Service) ListEntries(ctx context.Context, docset string) ([]*Entry, error) {
	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, fmt.Errorf("could not list entries in docset %q: %w", docset, err)
	}

	return idx.Entries(), nil
}

func (s *Service) ShowEntry(ctx context.Context, docset string, path string) (*EntryView, error) {
	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, fmt.Errorf("could not show entry %q in docset %q: %w", path, docset, err)
	}

	entry, ok := idx.Get(path)
	if !ok {
		return nil, fmt.Errorf("no entry %q found in docset %q", path, docset)
	}

	loc := NewEntryLocator(entry.Path)
	html, err := s.client.GetDocument(ctx, docset, loc)
	if err != nil {
		return nil, fmt.Errorf("could not fetch document for entry %q: %w", path, err)
	}

	md, err := s.converter.Convert(html)
	if err != nil {
		return nil, fmt.Errorf("could not convert entry %q to Markdown: %w", path, err)
	}

	var view *EntryView
	if loc.HasFragment() {
		lines, ok := md.Index.Get(loc.Fragment)
		if !ok {
			return nil, fmt.Errorf("searched for section %q in document %q: section not found", loc.Fragment, loc.Path)
		}

		view = NewExcerptView(md, lines)
	} else {
		view = NewDocumentView(md)
	}

	return view, err
}

func (s *Service) entryIndex(ctx context.Context, docset string) (*EntryIndex, error) {
	m, err := s.client.ListEntries(ctx, docset)
	if err != nil {
		return nil, fmt.Errorf("could not index entries in docset %q: %w", docset, err)
	}

	return NewEntryIndex(m.Entries), nil
}
