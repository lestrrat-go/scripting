package filter

import "io"

type Filter interface {
	Apply(io.Writer, io.Reader) error
}
