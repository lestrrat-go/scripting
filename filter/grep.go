package filter

import (
	"bufio"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

type Filter interface {
	Apply(io.Writer, io.Reader) error
}

type grep struct {
	pattern string
}

func Grep(pattern string) Filter {
	return &grep{pattern: pattern}
}

func (g *grep) Apply(dst io.Writer, src io.Reader) error {
	re, err := regexp.Compile(g.pattern)
	if err != nil {
		return errors.Wrapf(err, `failed to compile pattern '%s'`, g.pattern)
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if txt := scanner.Text(); re.MatchString(txt) {
			io.WriteString(dst, txt)
		}
	}
	return nil
}
