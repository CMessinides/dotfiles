package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-shellwords"
)

var (
	ErrNoDevDocsPager = errors.New("no DEVDOCS_PAGER environment variable set")
	ErrNoEnvPager     = errors.New("no PAGER environment variable set")
)

type PagerVars struct {
	Filename string
	Language string
}

type PagerOpts struct {
	Normalize bool
}

type Pager struct {
	Bin  string
	Args []string
	Env  map[string]string
}

func NewPager(cmd string, vars PagerVars, opts PagerOpts) (*Pager, error) {
	env := map[string]string{
		"DEVDOCS_FILENAME": vars.Filename,
		"DEVDOCS_LANGUAGE": vars.Language,
	}

	words, err := shellwords.Parse(expandWithCustomEnv(cmd, env))
	if err != nil {
		return nil, fmt.Errorf("could not parse pager command: %w", err)
	}

	bin := words[0]
	args := words[1:]

	if opts.Normalize {
		switch bin {
		case "less":
			args = append(args, "-+R", "-+F")
		}
	}

	bin, err = exec.LookPath(bin)
	if err != nil && !errors.Is(err, exec.ErrDot) {
		return nil, fmt.Errorf("could not determine path to pager: %w", err)
	}

	return &Pager{
		Bin:  bin,
		Args: args,
		Env:  env,
	}, nil
}

func expandWithCustomEnv(s string, vars map[string]string) string {
	return os.Expand(s, func(key string) string {
		if val, ok := vars[key]; ok {
			return val
		}

		return os.Getenv(key)
	})
}

func (p *Pager) Command() *exec.Cmd {
	cmd := exec.Command(p.Bin, p.Args...)

	cmd.Env = os.Environ()
	for k, v := range p.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	return cmd
}

func (p *Pager) Wrap(stdout io.WriteCloser, stderr io.WriteCloser) (*PagerWriter, error) {
	cmd := p.Command()
	slog.Debug("pager command", "cmd", cmd.String())

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not open pipe to pager: %w", err)
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start pager: %w", err)
	}

	return &PagerWriter{
		wc:  stdin,
		cmd: cmd,
	}, nil
}

func LookupPager(vars PagerVars) (*Pager, error) {
	p, err := lookupDevDocsPager(vars)

	if errors.Is(err, ErrNoDevDocsPager) {
		p, err = lookupEnvPager(vars)
	}

	if errors.Is(err, ErrNoEnvPager) {
		p, err = lookupDefaultPager(vars)
	}

	return p, err
}

func lookupDevDocsPager(vars PagerVars) (*Pager, error) {
	cmd, isSet := os.LookupEnv("DEVDOCS_PAGER")
	if !isSet || cmd == "" {
		return nil, ErrNoDevDocsPager
	}

	return NewPager(cmd, vars, PagerOpts{
		Normalize: false,
	})
}

func lookupEnvPager(vars PagerVars) (*Pager, error) {
	cmd, isSet := os.LookupEnv("PAGER")
	if !isSet || cmd == "" {
		return nil, ErrNoEnvPager
	}

	return NewPager(cmd, vars, PagerOpts{
		Normalize: true,
	})
}

func lookupDefaultPager(vars PagerVars) (*Pager, error) {
	return NewPager("less -R -F", vars, PagerOpts{
		Normalize: false,
	})
}

type PagerWriter struct {
	wc  io.WriteCloser
	cmd *exec.Cmd
}

func (w *PagerWriter) Write(p []byte) (n int, err error) {
	n, err = w.wc.Write(p)

	// Ignore pipe errors from closing the pager.
	if err != nil && errors.Is(err, syscall.EPIPE) {
		return len(p), nil
	}

	return n, err
}

func (w *PagerWriter) Close() error {
	errClose := w.wc.Close()
	errWait := w.cmd.Wait()
	return errors.Join(errClose, errWait)
}
