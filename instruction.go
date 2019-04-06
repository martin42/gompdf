package gompdf

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

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
	instructionRegistry.Register(&Image{})
}

type Instruction interface {
	Create(styles Styles, classes []string)
	ApplyClasses(scs StyleClasses)
}

type Instructions []Instruction

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
	allStyles := Styles{}
	allClasses := []string{}
	for _, a := range start.Attr {
		if a.Name.Local == "style" {
			styles, err := ParseStyles([]byte(a.Value))
			if err != nil {
				return nil, errors.Wrapf(err, "parse styles of element <%s> (%s)", start.Name.Local, a.Value)
			}
			allStyles = append(allStyles, styles...)
		} else if a.Name.Local == "class" {
			allClasses = append(allClasses, strings.Fields(a.Value)...)
		}
	}
	inst.Create(allStyles, allClasses)
	return inst, nil
}

type Styled struct {
	styles  Styles
	classes []string
}

func (i *Styled) Create(styles Styles, classes []string) {
	i.styles = styles
	i.classes = classes
}

func (i *Styled) ApplyClasses(scs StyleClasses) {
	for _, c := range i.classes {
		sc, found := scs.Find(c)
		if !found {
			continue
		}
		for _, s := range sc.Styles {
			if _, contains := i.styles.FindByPrototype(s); !contains {
				i.styles = append(i.styles, s)
			}
		}
	}
}

type NoStyles struct{}

func (i *NoStyles) Create(styles Styles, classes []string) {}

func (i *NoStyles) ApplyClasses(scs StyleClasses) {}

func (i *NoStyles) ApplyStyles() {}

type Font struct {
	Styled
	XMLName        xml.Name `xml:"Font"`
	fontFamily     FontFamily
	fontPointSize  FontPointSize
	fontStyle      FontStyle
	fontWeight     FontWeight
	fontDecoration FontDecoration
}

func (fnt *Font) ApplyStyles(def Font) {
	fnt.fontFamily = def.fontFamily
	fnt.fontPointSize = def.fontPointSize
	fnt.fontStyle = def.fontStyle
	for _, s := range fnt.styles {
		switch s := s.(type) {
		case FontFamily:
			fnt.fontFamily = s
		case FontStyle:
			fnt.fontStyle = s
		case FontPointSize:
			fnt.fontPointSize = s
		case FontWeight:
			fnt.fontWeight = s
		case FontDecoration:
			fnt.fontDecoration = s
		}
	}
}

type LineFeed struct {
	NoStyles
	XMLName xml.Name `xml:"Lf"`
	Lines   float64  `xml:"lines,attr"`
}

type SetX struct {
	NoStyles
	XMLName xml.Name `xml:"SetX"`
	Value   int      `xml:"value,attr"`
}

type SetY struct {
	NoStyles
	XMLName xml.Name `xml:"SetY"`
	Value   int      `xml:"value,attr"`
}

type Box struct {
	Styled
	XMLName xml.Name `xml:"Box"`
	Text    string   `xml:",chardata"`
}

func (b *Box) ApplyStyles() {

}

type Text struct {
	Styled
	XMLName    xml.Name `xml:"Text"`
	Text       string   `xml:",chardata"`
	lineHeight LineHeight
	width      Width
	hAlign     HAlign
}

func (t *Text) ApplyStyles(def Text) {
	t.lineHeight = def.lineHeight
	t.width = def.width
	t.hAlign = def.hAlign
	for _, s := range t.styles {
		switch s := s.(type) {
		case LineHeight:
			t.lineHeight = s
		case Width:
			t.width = s
		case HAlign:
			t.hAlign = s
		}
	}
}

type Image struct {
	Styled
	XMLName xml.Name `xml:"Image"`
	Source  string   `xml:",chardata"`
}

func (i *Image) ApplyStyles() {

}
