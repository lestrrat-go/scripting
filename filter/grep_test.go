package filter_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/lestrrat-go/scripting/filter"
	"github.com/stretchr/testify/assert"
)

func TestGrepBadPattern(t *testing.T) {
	g := filter.Grep(`(unterminated`)
	err := g.Apply(ioutil.Discard, &bytes.Buffer{})
	if !assert.Error(t, err, "g.Apply should fail") {
		return
	}

	err2 := g.Apply(ioutil.Discard, &bytes.Buffer{})
	if !assert.Equal(t, err, err2, "g.Apply should return same error") {
		return
	}
}
