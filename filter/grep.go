package filter

import (
	"bufio"
	"io"
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

type grep struct {
	pattern      string
	compiled     *regexp.Regexp
	compileError error
	compileOnce  sync.Once
}

func Grep(pattern string) Filter {
	return &grep{pattern: pattern}
}

func (g *grep) compilePattern() {
	re, err := regexp.Compile(g.pattern)
	if err != nil {
		g.compileError = errors.Wrapf(err, `failed to compile pattern '%s'`, g.pattern)
	} else {
		g.compiled = re
	}
}

func (g *grep) Apply(dst io.Writer, src io.Reader) error {
	g.compileOnce.Do(g.compilePattern)
	if g.compileError != nil {
		return g.compileError
	}

	re := g.compiled
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if txt := scanner.Text(); re.MatchString(txt) {
			io.WriteString(dst, txt)
		}
	}
	return nil
}
