package gompdf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	"github.com/mazzegi/gompdf/style"
	"github.com/pkg/errors"
)

var instructionRegistry *Registry

func init() {
	instructionRegistry = NewRegistry()
	instructionRegistry.Register(&Font{})
	instructionRegistry.Register(&Box{})
	instructionRegistry.Register(&Text{})
	instructionRegistry.Register(&LineFeed{})
	instructionRegistry.Register(&SetX{})
	instructionRegistry.Register(&SetY{})
	instructionRegistry.Register(&SetXY{})
	instructionRegistry.Register(&Image{})
	instructionRegistry.Register(&Table{})
	instructionRegistry.Register(&TableRow{})
	instructionRegistry.Register(&TableCell{})
}

type Instruction interface {
	DecodeAttrs(attrs []xml.Attr) error
	Apply(cs style.Classes, styles *style.Styles)
	ApplyWithSelector(sel string, cs style.Classes, styles *style.Styles)
}

type Instructions struct {
	Styled
	iss []Instruction
}

type Registry struct {
	types map[string]Instruction
}

func NewRegistry() *Registry {
	return &Registry{
		types: map[string]Instruction{},
	}
}

func (r *Registry) Register(prototype Instruction) error {
	ty := reflect.TypeOf(prototype)
	if ty.Kind() != reflect.Ptr {
		return errors.Errorf("register (%T). Instruction must be a ptr type (kind is %s).", ty.Name(), ty.Kind())
	}
	ty = reflect.TypeOf(reflect.ValueOf(prototype).Elem().Interface())
	fxml, ok := ty.FieldByName("XMLName")
	if !ok {
		return errors.Errorf("(%T) contains no XMLName", ty.Name())
	}
	xmlName := fxml.Tag.Get("xml")
	r.types[xmlName] = prototype
	return nil
}

func (r *Registry) Decode(d *xml.Decoder, start xml.StartElement) (Instruction, error) {
	proto, contains := r.types[start.Name.Local]
	if !contains {
		return nil, fmt.Errorf("registry-decode: (%s) is not registered", start.Name.Local)
	}
	pointerToI := reflect.New(reflect.TypeOf(proto))
	err := d.DecodeElement(pointerToI.Interface(), &start)
	if err != nil {
		return nil, err
	}
	inst := pointerToI.Elem().Interface().(Instruction)
	err = inst.DecodeAttrs(start.Attr)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

type Styled struct {
	Appliers []*style.Applier
	Classes  []string
}

func (i *Styled) DecodeAttrs(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Local == "style" {
			app, err := style.DecodeApplier(bytes.NewBufferString(a.Value))
			if err != nil {
				return errors.Wrapf(err, "decode style applier (%s)", a.Value)
			}
			i.Appliers = append(i.Appliers, app)
		} else if a.Name.Local == "class" {
			i.Classes = append(i.Classes, strings.Fields(a.Value)...)
		}
	}
	return nil
}

func (i *Styled) Apply(cs style.Classes, styles *style.Styles) {
	cs.Apply(styles, i.Classes...)
	for _, app := range i.Appliers {
		app.Apply(styles)
	}
}

func (i *Styled) ApplyWithSelector(sel string, cs style.Classes, styles *style.Styles) {
	cs.ApplyWithSelector(sel, styles, i.Classes...)
	for _, app := range i.Appliers {
		app.Apply(styles)
	}
}

type NoStyles struct{}

func (i *NoStyles) DecodeAttrs(attrs []xml.Attr) error { return nil }

func (i *NoStyles) Apply(cs style.Classes, styles *style.Styles) {}

func (i *NoStyles) ApplyWithSelector(sel string, cs style.Classes, styles *style.Styles) {}

type Font struct {
	Styled
	XMLName xml.Name `xml:"Font"`
}

type LineFeed struct {
	NoStyles
	XMLName xml.Name `xml:"Lf"`
	Lines   float64  `xml:"lines,attr"`
}

type SetX struct {
	NoStyles
	XMLName xml.Name `xml:"SetX"`
	X       float64  `xml:"x,attr"`
}

type SetY struct {
	NoStyles
	XMLName xml.Name `xml:"SetY"`
	Y       float64  `xml:"y,attr"`
}

type SetXY struct {
	NoStyles
	XMLName xml.Name `xml:"SetXY"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
}

type Box struct {
	Styled
	XMLName xml.Name `xml:"Box"`
	Text    string   `xml:",chardata"`
}

type Text struct {
	Styled
	XMLName xml.Name `xml:"Text"`
	Text    string   `xml:",chardata"`
}

type Image struct {
	Styled
	XMLName xml.Name `xml:"Image"`
	Source  string   `xml:",chardata"`
}
