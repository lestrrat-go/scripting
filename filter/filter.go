package filter

import (
	"bufio"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

func Grep(dst io.Writer, src io.Reader, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, `failed to compile pattern '%s'`, pattern)
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		if txt := scanner.Text(); re.MatchString(txt) {
			io.WriteString(dst, txt)
		}
	}
	return nil
}
