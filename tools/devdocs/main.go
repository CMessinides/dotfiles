package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"golang.org/x/term"
)

type Context struct {
	context.Context
	Renderer Renderer
	Service  *Service
}

type DocsetsListCmd struct{}

func (c DocsetsListCmd) Run(ctx *Context) error {
	docsets, err := ctx.Service.ListDocsets(ctx)
	if err != nil {
		return err
	}

	return ctx.Renderer.RenderDocsetList(docsets)
}

type EntriesListCmd struct {
	Docset string `arg:"" help:"Docset to retrieve"`
}

func (c EntriesListCmd) Run(ctx *Context) error {
	entries, err := ctx.Service.ListEntries(ctx, c.Docset)
	if err != nil {
		return err
	}

	return ctx.Renderer.RenderEntryList(entries)
}

type EntriesShowCmd struct {
	Docset string `arg:"" help:"Docset to retrieve documentation from"`
	Path   string `arg:"" help:"Path to the entry"`
}

func (c EntriesShowCmd) Run(ctx *Context) error {
	view, err := ctx.Service.ShowEntry(ctx, c.Docset, c.Path)
	if err != nil {
		return err
	}

	return ctx.Renderer.RenderEntryView(view)
}

type CLI struct {
	Debug     bool   `help:"Enable debug mode"`
	Format    string `help:"Specify the output format" default:"console" enum:"console,porcelain,json"`
	JSON      bool   `help:"Print JSON. Shortcut for --format=json" xor:"fmt"`
	Porcelain bool   `help:"Print script-friendly text. Shortcut for --format=porcelain" xor:"fmt"`
	Docsets   struct {
		List DocsetsListCmd `cmd:"" help:"List all docsets"`
	} `cmd:"" help:"Get information about docsets"`

	Entries struct {
		List EntriesListCmd `cmd:"" help:"List all entries in a docset"`
		Show EntriesShowCmd `cmd:"" help:"Show documentation for an entry"`
	} `cmd:"" help:"Get information about entries"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(
		&cli,
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	if cli.Debug {
		SetLogLevel(slog.LevelDebug)
	}

	var renderer Renderer
	// Prefer the shortcut flags --json and --porcelain.
	if cli.JSON {
		cli.Format = "json"
	} else if cli.Porcelain {
		cli.Format = "porcelain"
	}
	switch cli.Format {
	case "json":
		renderer = NewJSONRenderer(os.Stdout)
	case "porcelain":
		renderer = NewPorcelainRenderer(os.Stdout)
	default:
		isTTY := term.IsTerminal(int(os.Stderr.Fd()))
		renderer = NewConsoleRenderer(os.Stdout, os.Stderr, isTTY)
	}

	err := ctx.Run(&Context{
		Context:  context.Background(),
		Renderer: renderer,
		Service: NewService(
			DefaultCache,
			DefaultClient,
			DefaultMarkdownConverter,
		),
	})
	ctx.FatalIfErrorf(err)
}
