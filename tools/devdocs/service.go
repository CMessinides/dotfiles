package main

import (
	"context"
	"fmt"

	"github.com/cmessinides/dotfiles/tools/devdocs/internal/report"
)

type EntryNotFoundError struct {
	report.Err
	Docset string
	Path   string
}

type SectionNotFoundError struct {
	report.Err
	ID     string
	Path   string
	Docset string
}

type Service struct {
	cache     Cache
	client    *Client
	converter *MarkdownConverter
	errs      *report.Chain
}

func NewService(cache Cache, client *Client, converter *MarkdownConverter) *Service {
	return &Service{
		cache:     cache,
		client:    client,
		converter: converter,
		errs: report.NewChain(
			report.WithPrefix("Service"),
		),
	}
}

func (s *Service) ListDocsets(ctx context.Context) ([]Docset, error) {
	errs := s.errs.Extend(
		report.WithMethodLabel("ListDocsets"),
	)

	d, err := s.client.ListDocsets(ctx)
	if err != nil {
		return d, errs.Wrap("could not get docsets", err)
	}

	return d, nil
}

func (s *Service) ListEntries(ctx context.Context, docset string) ([]*Entry, error) {
	errs := s.errs.Extend(
		report.WithMethodLabel("ListEntries", docset),
	)

	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, errs.Wrap("could not index entries", err)
	}

	return idx.Entries(), nil
}

func (s *Service) ShowEntry(ctx context.Context, docset string, path string) (*EntryView, error) {
	errs := s.errs.Extend(
		report.WithMethodLabel("ShowEntry", docset, path),
	)

	idx, err := s.entryIndex(ctx, docset)
	if err != nil {
		return nil, errs.Wrap("could not index entries", err)
	}

	entry, ok := idx.Get(path)
	if !ok {
		return nil, &EntryNotFoundError{
			Err: errs.New(
				fmt.Sprintf(
					"docset %q index has no entry %q", docset, path),
			),
			Docset: docset,
			Path:   path,
		}
	}

	loc := NewEntryLocator(entry.Path)
	html, err := s.client.GetDocument(ctx, docset, loc)
	if err != nil {
		return nil, errs.Wrap("could not get document", err)
	}

	md, err := s.converter.Convert(html)
	if err != nil {
		return nil, errs.Wrap("could not convert document to Markdown", err)
	}

	var view *EntryView
	if loc.HasFragment() {
		lines, ok := md.Index.Get(loc.Fragment)
		if !ok {
			return nil, &SectionNotFoundError{
				Err: errs.New(
					fmt.Sprintf("document %q has no section with ID %q", loc.Path, loc.Fragment),
				),
				ID:     loc.Fragment,
				Path:   loc.Path,
				Docset: docset,
			}
		}

		view = NewExcerptView(md, lines)
	} else {
		view = NewDocumentView(md)
	}

	return view, nil
}

func (s *Service) entryIndex(ctx context.Context, docset string) (*EntryIndex, error) {
	m, err := s.client.ListEntries(ctx, docset)
	if err != nil {
		return nil, err
	}

	return NewEntryIndex(m.Entries), nil
}
