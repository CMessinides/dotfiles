package main

import (
	"bytes"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

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

type MarkdownConverter struct {
	Preprocessors []HTMLPreprocessor
}

func (m *MarkdownConverter) Convert(src *HTMLDocument) (*MarkdownDocument, error) {
	html, err := goquery.NewDocumentFromReader(src.Content.Reader())
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	sel := html.Selection
	for _, p := range m.Preprocessors {
		sel, err = p.Preprocess(sel)
		if err != nil {
			return nil, fmt.Errorf("failed to preprocess HTML: %w", err)
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
			return nil, fmt.Errorf("failed to convert node to Markdown: %w", err)
		}

		buf.Write(md)
	}

	data := buf.Bytes()

	idx, err := BuildDocumentIndex(data, ids)
	if err != nil {
		return nil, err
	}

	return NewMarkdownDocumentFromHTML(src, data, idx), nil
}

type MarkdownConverterConfigFunc func(m *MarkdownConverter)

func WithPreprocessors(p ...HTMLPreprocessor) MarkdownConverterConfigFunc {
	return func(m *MarkdownConverter) {
		m.Preprocessors = append(m.Preprocessors, p...)
	}
}

func NewMarkdownConverter(configs ...MarkdownConverterConfigFunc) *MarkdownConverter {
	m := &MarkdownConverter{
		Preprocessors: make([]HTMLPreprocessor, 0),
	}

	for _, configure := range configs {
		configure(m)
	}

	return m
}

var DefaultMarkdownConverter = NewMarkdownConverter(
	WithPreprocessors(
		NormalizeLanguagesOnCodeBlocks,
		AddLanguageClassesToCodeBlocks,
	),
)
