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

var ErrNotFound = errors.New("resource not found")

type DocsetList struct {
	Docsets []Docset `json:"docsets"`
}

type Client struct {
	*http.Client
	rootURL      *url.URL
	documentsURL *url.URL
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
	}
}

func mustParseURL(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return url
}

func (c *Client) ListDocsets(ctx context.Context) (DocsetList, error) {
	list := DocsetList{
		Docsets: make([]Docset, 0),
	}

	u := c.rootURL.JoinPath("/docs/docs.json").String()
	res, err := c.get(ctx, u)
	if err != nil {
		return list, err
	}

	j := json.NewDecoder(res.Body)
	err = j.Decode(&list.Docsets)
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

	u := c.rootURL.JoinPath("/docs/", docset, "/index.json").String()
	res, err := c.get(ctx, u)
	if err != nil {
		return m, fmt.Errorf("searched for docset with slug \"%s\": %w", docset, err)
	}

	j := json.NewDecoder(res.Body)
	err = j.Decode(&m)
	if err != nil {
		return m, err
	}

	err = res.Body.Close()
	if err != nil {
		return m, nil
	}

	m.Docset = docset
	return m, nil
}

func (c *Client) GetDocument(ctx context.Context, docset string, path string, fragment string) (Document, error) {
	u := c.documentsURL.JoinPath("/", docset, "/", path+".html").String()
	res, err := c.get(ctx, u)
	if err != nil {
		return EmptyDocument, fmt.Errorf("searched for path \"%s\" in docset \"%s\": %w", path, docset, err)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return EmptyDocument, err
	}

	err = res.Body.Close()
	if err != nil {
		return EmptyDocument, err
	}

	return NewDocument(docset, path, fragment, ContentHTML, buf.Bytes()), nil
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	slog.Debug("initiating request", "url", url, "method", "GET")
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		slog.Debug("failed to create request", "url", url, "err", err)
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		slog.Debug("failed to get response", "url", url, "err", err)
		return nil, err
	}

	slog.Debug(
		"got response from DevDocs",
		"url", url,
		"status", res.Status,
		"content-type", res.Header.Get("Content-Type"),
	)

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("received HTTP error from DevDocs: %s", res.Status)
		}
	}

	return res, nil
}

var DefaultClient = NewClient(ClientOptions{
	Client:       httpClient,
	RootURL:      DefaultDevDocsURL,
	DocumentsURL: DefaultDevDocsDocumentsURL,
})
