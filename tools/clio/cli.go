package clio

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
)

type cliContext struct{}

type addCmd struct {
	Path string `arg:"" type:"existingfile" help:"The script to add."`
}

func (a *addCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from add")
	return nil
}

type defaultCmd struct{}

func (s *defaultCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from default")
	return nil
}

type filterCmd struct {
	Args []string `arg:"" help:"Args to pass through to fzf."`
}

func (f *filterCmd) Run(ctx *cliContext) error {
	fzf := exec.Command("fzf", f.Args...)
	fzf.Stdin = os.Stdin
	fzf.Stdout = os.Stdout
	fzf.Stderr = os.Stderr

	err := fzf.Run()
	return err
}

type helpCmd struct {
	Script     string `arg:"" `
	Subcommand string `arg:"" optional:"" help:"Narrow help "`
}

func (h *helpCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from help")
	return nil
}

type listCmd struct{}

func (l *listCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from list")
	return nil
}

type removeCmd struct {
	Script string `arg:""`
}

func (r *removeCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from remove")
	return nil
}

type runCmd struct {
	Script string   `arg:""`
	Args   []string `arg:"" passthrough:""`
}

func (r *runCmd) Run(ctx *cliContext) error {
	fmt.Println("hello from run")
	return nil
}

var cli struct {
	Add     addCmd     `cmd:"" help:"Add a script."`
	Default defaultCmd `cmd:"" hidden:"" default:"1"`
	Filter  filterCmd  `cmd:"" help:"Filter standard input through fzf." passthrough:""`
	Help    helpCmd    `cmd:"" help:"Display help for a script."`
	List    listCmd    `cmd:"" help:"List available scripts."`
	Remove  removeCmd  `cmd:"" help:"Remove a script."`
	Run     runCmd     `cmd:"" help:"Run a script."`
}

func Run() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&cliContext{})
	ctx.FatalIfErrorf(err)
}
