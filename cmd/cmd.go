package cmd

import (
	"bytes"
	"context"
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
	grep          string
	path          string
	spinner       bool
}

type Result struct {
	output *bytes.Buffer
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

func New(path string, args ...string) *Command {
	return &Command{
		path: path,
		args: args,
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

// Grep specifies that filtering on the output is performed
// after the command has been executed. `pattern` is treated
// as a regular expression. DO NOT forget to call one or both
// of `CaptureStderr` or `CaptureStdout`, otherwise there will
// be nothing to filter against.
func (c *Command) Grep(pattern string) *Command {
	c.grep = pattern
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
	ec.grep = c.grep
	ec.path = c.path
	ec.spinner = c.spinner

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
			return nil, errors.Wrap(err, `failed to execute command`)
		}
	}
	close(done)

	if out != nil {
		if pattern := c.grep; len(c.grep) > 0 {
			var dst bytes.Buffer
			if err := filter.Grep(&dst, out, pattern); err != nil {
				return nil, errors.Wrap(err, `failed to apply grep`)
			}
			out = &dst
		}
	}

	return &Result{
		output: out,
	}, nil
}
