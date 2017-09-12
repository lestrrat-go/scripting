package cmd_test

import (
	"testing"

	"github.com/lestrrat/go-scripting/cmd"
	"github.com/stretchr/testify/assert"
)

func TestSed(t *testing.T) {
	r, err := cmd.New("cat", "testdata/example.txt").
		CaptureStderr(true).
		CaptureStdout(true).
		Sed("hello", "こんにちわ").
		Do(nil)
	if !assert.NoError(t, err, "command should succeed") {
		return
	}

	want := "こんにちわ foo\nbar こんにちわ\nherro baz\n"
	if !assert.Equal(t, want, r.OutputString(), "want %v, got %v", want, r.OutputString()) {
		return
	}
}
