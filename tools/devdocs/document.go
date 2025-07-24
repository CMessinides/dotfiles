package main

import (
	"bytes"
	"fmt"
	"log/slog"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

type DocumentTransformer interface {
	Transform(doc Document) (Document, error)
}

type HTMLPreprocessor interface {
	Preprocess(s *goquery.Selection) (*goquery.Selection, error)
}

type HTMLPreprocessorFunc func(s *goquery.Selection) (*goquery.Selection, error)

func (h HTMLPreprocessorFunc) Preprocess(s *goquery.Selection) (*goquery.Selection, error) {
	return h(s)
}

var NormalizeLanguagesOnCodeBlocks = HTMLPreprocessorFunc(func(s *goquery.Selection) (*goquery.Selection, error) {
	s.Find("[data-language]").Each(func(i int, s *goquery.Selection) {
		lang, _ := s.Attr("data-language")
		switch lang {
		case "go":
			// Compatibility with bat syntax highlighting.
			s.SetAttr("data-language", "golang")
		}
	})

	return s, nil
})

var AddLanguageClassesToCodeBlocks = HTMLPreprocessorFunc(func(s *goquery.Selection) (*goquery.Selection, error) {
	s.Find("[data-language]").Each(func(i int, s *goquery.Selection) {
		lang, _ := s.Attr("data-language")
		s.AddClass("language-" + lang)
	})

	return s, nil
})

type MarkdownTransformer struct {
	Preprocessors []HTMLPreprocessor
}

func (m *MarkdownTransformer) Transform(doc Document) (Document, error) {
	html, err := goquery.NewDocumentFromReader(
		bytes.NewReader(doc.Content),
	)
	if err != nil {
		return EmptyDocument, fmt.Errorf("failed to parse HTML: %w", err)
	}

	sel := html.Selection
	for _, p := range m.Preprocessors {
		sel, err = p.Preprocess(sel)
		if err != nil {
			return EmptyDocument, fmt.Errorf("failed to preprocess HTML: %w", err)
		}
	}

	headings := sel.Find("h2, h3, h4, h5, h6")
	ids := headings.Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("id", "")
	})

	buf := new(bytes.Buffer)
	for _, node := range sel.Nodes {
		md, err := htmltomarkdown.ConvertNode(node)
		if err != nil {
			return EmptyDocument, fmt.Errorf("failed to convert node to Markdown: %w", err)
		}

		buf.Write(md)
	}

	data := buf.Bytes()

	toc, err := BuildTableOfContents(data, ids)
	if err != nil {
		return EmptyDocument, err
	}

	slog.Debug("built table of contents", "len", len(toc.Sections))

	return doc.ToMarkdown(data), nil
}

type MarkdownTransformerConfigFunc func(m *MarkdownTransformer)

func WithPreprocessors(p ...HTMLPreprocessor) MarkdownTransformerConfigFunc {
	return func(m *MarkdownTransformer) {
		m.Preprocessors = append(m.Preprocessors, p...)
	}
}

func NewMarkdownTransformer(configs ...MarkdownTransformerConfigFunc) *MarkdownTransformer {
	m := &MarkdownTransformer{
		Preprocessors: make([]HTMLPreprocessor, 0),
	}

	for _, configure := range configs {
		configure(m)
	}

	return m
}

var DefaultMarkdownTransformer = NewMarkdownTransformer(
	WithPreprocessors(
		NormalizeLanguagesOnCodeBlocks,
		AddLanguageClassesToCodeBlocks,
	),
)
