package cmd_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/lestrrat/go-scripting/cmd"
	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	res, err := cmd.New("ls", "-l").
		CaptureStdout(true).
		Grep(`_test\.go$`).
		Do(nil)
	if !assert.NoError(t, err, "ls -l should succeed") {
		return
	}

	// Each line should ONLY contain lines with _test.go$
	scanner := bufio.NewScanner(res.Output())
	for scanner.Scan() {
		txt := scanner.Text()
		t.Logf("got '%s'", txt)
		if !assert.True(t, strings.HasSuffix(txt, "_test.go"), "each line should contain _test.go") {
			return
		}
	}
}

func ExampleCommand() {
	_, err := cmd.New("ls", "-l").
		CaptureStdout(true).
		Grep(`_test\.go$`).
		Do(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// OUTPUT:
}
