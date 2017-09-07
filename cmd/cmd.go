package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/lestrrat/go-scripting/filter"
	"github.com/pkg/errors"
)

type Command struct {
	args          []string
	bailOnErr     bool
	captureStdout bool
	captureStderr bool
	filters       []filter.Filter
	path          string
	spinner       bool
	stdin         io.Reader
}

type Result struct {
	output *bytes.Buffer
}

// JSON assumes that your accumulated output is a JSON string, and
// attemps to decode it.
func (r *Result) JSON(v interface{}) error {
	return json.NewDecoder(r.output).Decode(v)
}

func (r *Result) Output() *bytes.Buffer {
	return r.output
}

func (r *Result) OutputString() string {
	if r.output == nil {
		return ""
	}

	return r.output.String()
}

func Exec(path string, args ...string) error {
	_, err := New(path, args...).BailOnError(true).Do(nil)
	return err
}

// New creates a new Command instance. `path` (i.e. the command to
// execute) is required. By default, `BailOnError` is true
func New(path string, args ...string) *Command {
	return &Command{
		bailOnErr: true,
		path:      path,
		args:      args,
	}
}

func (c *Command) BailOnError(b bool) *Command {
	c.bailOnErr = b
	return c
}

func (c *Command) CaptureStderr(b bool) *Command {
	c.captureStderr = b
	return c
}

func (c *Command) CaptureStdout(b bool) *Command {
	c.captureStdout = b
	return c
}

func (c *Command) Spinner(b bool) *Command {
	c.spinner = b
	return c
}

// Stdin sets a stdin to be piped to the command
func (c *Command) Stdin(in io.Reader) *Command {
	c.stdin = in
	return c
}

// Grep adds a new filtering on the output of the command
// after it has been executed. `pattern` is treated
// as a regular expression. DO NOT forget to call one or both
// of `CaptureStderr` or `CaptureStdout`, otherwise there will
// be nothing to filter against.
func (c *Command) Grep(pattern string) *Command {
	return c.Filter(filter.Grep(pattern))
}

func (c *Command) Sed(pattern, replace string) *Command {
	return c.Filter(filter.Sed(pattern, replace))
}

func (c *Command) Filter(f filter.Filter) *Command {
	c.filters = append(c.filters, f)
	return c
}

// Do executes the command, and applies other filters as necessary
func (c *Command) Do(ctx context.Context) (*Result, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}

	// create an execution context, so we don't have to worry about
	// people mutating `c` after Do() is called, or if they call Do()
	// multiple times and stuff
	var ec execCtx
	ec.args = c.args
	ec.bailOnErr = c.bailOnErr
	ec.captureStdout = c.captureStdout
	ec.captureStderr = c.captureStderr
	ec.filters = c.filters
	ec.path = c.path
	ec.spinner = c.spinner
	ec.stdin = c.stdin

	return ec.Do(ctx)
}

type execCtx Command

func (c *execCtx) Do(ctx context.Context) (*Result, error) {
	done := make(chan struct{})

	cmd := exec.CommandContext(ctx, c.path, c.args...)
	var out *bytes.Buffer
	if c.captureStdout || c.captureStderr {
		out = &bytes.Buffer{}
		if c.captureStdout {
			cmd.Stdout = out
		}
		if c.captureStderr {
			cmd.Stderr = out
		}
	}

	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}

	// Start a spinner
	if c.spinner {
		// XXX for now, this is hardcoded
		s := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
		s.Start()
		go func() {
			defer s.Stop()
			select {
			case <-ctx.Done():
			case <-done:
			}
		}()
	}

	if err := cmd.Run(); err != nil {
		if c.bailOnErr {
			return &Result{output: out}, errors.Wrap(err, `failed to execute command`)
		}
	}
	close(done)

	if out != nil {
		for _, f := range c.filters {
			var dst bytes.Buffer
			if err := f.Apply(&dst, out); err != nil {
				return nil, errors.Wrapf(err, `failed to apply filter %s`, f)
			}
			out = &dst
		}
	}

	return &Result{
		output: out,
	}, nil
}
