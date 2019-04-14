package style

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Class struct {
	Name    string
	applier *Applier
}

type Classes map[string]Class

func (c Class) Apply(styles *Styles) {
	c.applier.Apply(styles)
}

func (cs Classes) Apply(styles *Styles, classes ...string) {
	for _, cn := range classes {
		if c, ok := cs[cn]; ok {
			c.Apply(styles)
		}
	}
}

func DecodeClasses(r io.Reader) (Classes, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read-all")
	}
	s := string(b)

	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	cs := Classes{}
	pos := 0
	for {
		curr := s[pos:]
		i := strings.IndexByte(curr, '{')
		if i < 0 {
			return cs, nil
		}
		name := trimWS(curr[:i])
		if len(name) == 0 {
			return nil, errors.Errorf("style class without name")
		}
		in := strings.IndexByte(curr[i:], '}')
		if in < 0 {
			return nil, errors.Errorf("non matching brace")
		}
		in += i
		app, err := DecodeApplier(bytes.NewBufferString(curr[i+1 : in]))
		if err != nil {
			return nil, errors.Wrap(err, "parse style")
		}
		cs[name] = Class{
			Name:    string(name),
			applier: app,
		}
		pos += in + 1
	}
}
