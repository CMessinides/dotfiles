package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultDevDocsURL          = "https://devdocs.io/"
	DefaultDevDocsDocumentsURL = "https://documents.devdocs.io"
)

type HTTPResponseError struct {
	LabeledError
	StatusCode int
	Status     string
}

func isNotFound(err error) bool {
	var h *HTTPResponseError
	if errors.As(err, &h) {
		return h.StatusCode == http.StatusNotFound
	}

	return false
}

type DocsetNotFoundError struct {
	LabeledError
	Docset string
}

// Error implements the error interface.
func (d *DocsetNotFoundError) Error() string {
	return fmt.Sprintf("docset %q not found", d.Docset)
}

type DocumentNotFoundError struct {
	LabeledError
	Docset string
	Entry  EntryLocator
}

// Error implements the error interface.
func (e *DocumentNotFoundError) Error() string {
	return fmt.Sprintf("document %q not found in docset %q", e.Entry, e.Docset)
}

var ErrNotFound = errors.New("resource not found")

type Client struct {
	*http.Client
	rootURL      *url.URL
	documentsURL *url.URL
	errs         *ErrorBuilder
}

type ClientOptions struct {
	Client       *http.Client
	RootURL      string
	DocumentsURL string
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func NewClient(opts ClientOptions) *Client {
	c := opts.Client
	if c == nil {
		c = httpClient
	}

	return &Client{
		Client:       c,
		rootURL:      mustParseURL(opts.RootURL),
		documentsURL: mustParseURL(opts.DocumentsURL),
		errs:         NewErrorBuilder(WithPrefix("client")),
	}
}

func mustParseURL(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return url
}

func (c *Client) ListDocsets(ctx context.Context) ([]Docset, error) {
	list := make([]Docset, 0)

	u := c.rootURL.JoinPath("/docs/docs.json").String()
	res, err := c.get(ctx, u)
	if err != nil {
		return list, err
	}

	j := json.NewDecoder(res.Body)
	err = j.Decode(&list)
	if err != nil {
		return list, err
	}

	err = res.Body.Close()
	if err != nil {
		return list, err
	}

	return list, nil
}

func (c *Client) ListEntries(ctx context.Context, docset string) (EntryManifest, error) {
	m := EntryManifest{
		Entries: make([]Entry, 0),
	}

	errs := c.errs.Extend(WithMethodLabel("ListEntries", docset))

	u := c.rootURL.JoinPath("/docs/", docset, "/index.json").String()
	res, err := c.get(ctx, u)
	if err != nil {
		if isNotFound(err) {
			return m, &DocsetNotFoundError{
				LabeledError: errs.Wrap(
					fmt.Sprintf("docset %q not found", docset),
					err,
				),
				Docset: docset,
			}
		} else {
			return m, errs.Wrap("could not get entry manifest", err)
		}
	}

	j := json.NewDecoder(res.Body)
	err = j.Decode(&m)
	if err != nil {
		return m, errs.Wrap("could not decode entry manifest JSON", err)
	}

	err = res.Body.Close()
	if err != nil {
		return m, errs.Wrap("could not close response stream", err)
	}

	m.Docset = docset
	return m, nil
}

func (c *Client) GetDocument(ctx context.Context, docset string, entry EntryLocator) (*HTMLDocument, error) {
	errs := c.errs.Extend(
		WithMethodLabel("GetDocument", docset, entry),
	)

	u := c.documentsURL.JoinPath("/", docset, "/", entry.Path+".html").String()
	res, err := c.get(ctx, u)
	if err != nil {
		if isNotFound(err) {
			return nil, &DocumentNotFoundError{
				LabeledError: errs.Wrap(
					fmt.Sprintf("document %q not found in docset %q", entry, docset),
					err,
				),
			}
		} else {
			return nil, c.errs.Wrap("could not get document", err)
		}
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, c.errs.Wrap("could not read response stream", err)
	}

	err = res.Body.Close()
	if err != nil {
		return nil, c.errs.Wrap("could not close response stream", err)
	}

	return NewHTMLDocument(docset, entry, buf.Bytes()), nil
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	slog.Debug("initiating request", "url", url, "method", "GET")
	errs := c.errs.Extend(
		WithMethodLabel("get", url),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		slog.Debug("failed to create request", "url", url, "err", err)
		return nil, errs.Wrap("failed to create request", err)
	}

	res, err := c.Do(req)
	if err != nil {
		slog.Debug("failed to get response", "url", url, "err", err)
		return nil, errs.Wrap("failed to get response", err)
	}

	slog.Debug(
		"got response from DevDocs",
		"url", url,
		"status", res.Status,
		"content-type", res.Header.Get("Content-Type"),
	)

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, &HTTPResponseError{
			LabeledError: errs.New(
				fmt.Sprintf("got HTTP response error from DevDocs: %s", res.Status),
			),
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}
	}

	return res, nil
}

var DefaultClient = NewClient(ClientOptions{
	Client:       httpClient,
	RootURL:      DefaultDevDocsURL,
	DocumentsURL: DefaultDevDocsDocumentsURL,
})
