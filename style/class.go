package style

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Selector struct {
	Name    string
	applier *Applier
}

type Class struct {
	Name      string
	applier   *Applier
	Selectors map[string]Selector
}

type Classes map[string]Class

func (c Class) Apply(styles *Styles) {
	c.applier.Apply(styles)
}

func (c Class) ApplyWithSelector(sel string, styles *Styles) {
	if sel, ok := c.Selectors[sel]; ok {
		sel.applier.Apply(styles)
		return
	}
	c.applier.Apply(styles)
}

func (cs Classes) Apply(styles *Styles, classes ...string) {
	for _, cn := range classes {
		if c, ok := cs[cn]; ok {
			c.Apply(styles)
		}
	}
}

func (cs Classes) ApplyWithSelector(sel string, styles *Styles, classes ...string) {
	for _, cn := range classes {
		if c, ok := cs[cn]; ok {
			c.ApplyWithSelector(sel, styles)
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
		className := name
		selName := ""
		idxSel := strings.Index(name, ":")
		if idxSel > 0 {
			className = name[:idxSel]
			selName = name[idxSel+1:]
			if cl, ok := cs[className]; ok {
				cl.Selectors[selName] = Selector{
					Name:    selName,
					applier: app,
				}
				fmt.Printf("added selector (%s) to class (%s) \n", selName, className)
			} else {
				return nil, errors.Errorf("no base class for (%s:%s)", className, selName)
			}
		} else {
			//no selector
			cl := Class{
				Name:      string(className),
				applier:   app,
				Selectors: map[string]Selector{},
			}
			cs[className] = cl
		}
		pos += in + 1
	}
}
