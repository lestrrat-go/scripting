package filter

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

type sed struct {
	pattern string
	replace string
}

func Sed(pattern, replace string) Filter {
	return &sed{pattern: pattern, replace: replace}
}

func (s *sed) Apply(dst io.Writer, src io.Reader) error {
	re, err := regexp.Compile(s.pattern)
	if err != nil {
		return errors.Wrapf(err, `failed to compile pattern '%s'`, s.pattern)
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		text := scanner.Text()
		if m := re.FindAllStringIndex(text, -1); len(m) > 0 {
			var buf bytes.Buffer
			buf.WriteString(text[0:m[0][0]])
			for i := range m {
				var e int
				if i == len(m)-1 {
					e = len(text)
				} else {
					e = m[i+1][0]
				}
				buf.WriteString(s.replace + text[m[i][1]:e])
			}
			io.WriteString(dst, buf.String()+"\n")
		} else {
			io.WriteString(dst, text+"\n")
		}
	}
	return nil
}
