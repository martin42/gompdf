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
	instructionRegistry.Register(&SetXY{})
	instructionRegistry.Register(&Image{})
	instructionRegistry.Register(&Table{})
}

type Instruction interface {
	Create(styles Styles, classes []string)
	ApplyClasses(scs StyleClasses)
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

	Logf("decoded: %s", start.Name.Local)
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
	Styles  Styles   //`xml:"style,attr"`
	Classes []string //`xml:"class,attr"`
}

// func (s *Styled) UnmarshalXMLAttr(attr xml.Attr) error {
// 	Logf("UnmarshalXMLAttr: %s", attr.Name.Local)
// 	if attr.Name.Local == "style" {
// 		styles, err := ParseStyles([]byte(attr.Value))
// 		if err != nil {
// 			return err
// 		}
// 		s.Styles = styles
// 	} else if attr.Name.Local == "class" {
// 		s.Classes = strings.Fields(attr.Value)
// 	}
// 	return nil
// }

func (i *Styled) Create(styles Styles, classes []string) {
	i.Styles = styles
	i.Classes = classes
}

func (i *Styled) ApplyClasses(scs StyleClasses) {
	for _, c := range i.Classes {
		sc, found := scs.Find(c)
		if !found {
			continue
		}
		for _, s := range sc.Styles {
			if _, contains := i.Styles.FindByPrototype(s); !contains {
				i.Styles = append(i.Styles, s)
			}
		}
	}
}

type NoStyles struct{}

func (i *NoStyles) Create(styles Styles, classes []string) {}

func (i *NoStyles) ApplyClasses(scs StyleClasses) {}

func (i *NoStyles) ApplyStyles() {}

type FontStyles struct {
	fontFamily     FontFamily
	fontPointSize  FontPointSize
	fontStyle      FontStyle
	fontWeight     FontWeight
	fontDecoration FontDecoration
}

func (fs *FontStyles) ApplyStyles(def FontStyles, styles Styles) {
	*fs = def
	for _, s := range styles {
		switch s := s.(type) {
		case FontFamily:
			fs.fontFamily = s
		case FontStyle:
			fs.fontStyle = s
		case FontPointSize:
			fs.fontPointSize = s
		case FontWeight:
			fs.fontWeight = s
		case FontDecoration:
			fs.fontDecoration = s
		}
	}
}

type Font struct {
	Styled
	FontStyles
	XMLName xml.Name `xml:"Font"`
}

func (fnt *Font) ApplyStyles(def FontStyles) {
	fnt.FontStyles.ApplyStyles(def, fnt.Styles)
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

type DrawingStyles struct {
	backgroundColor BackgroundColor
	color           Color
	lineWidth       LineWidth
}

func (b *DrawingStyles) ApplyStyles(def DrawingStyles, styles Styles) {
	*b = def
	for _, s := range styles {
		switch s := s.(type) {
		case BackgroundColor:
			b.backgroundColor = s
		case Color:
			b.color = s
		case LineWidth:
			b.lineWidth = s
		}
	}
}

type BoxStyles struct {
	border  Border
	padding Padding
}

func (b *BoxStyles) ApplyStyles(def BoxStyles, styles Styles) {
	*b = def
	for _, s := range styles {
		switch s := s.(type) {
		case Border:
			b.border = s
		case Padding:
			b.padding = s
		}
	}
}

type Box struct {
	Styled
	BoxStyles
	TextStyles
	DrawingStyles
	XMLName xml.Name `xml:"Box"`
	Text    string   `xml:",chardata"`
}

func (b *Box) ApplyStyles(defBox BoxStyles, defText TextStyles, defDraw DrawingStyles) {
	b.BoxStyles.ApplyStyles(defBox, b.Styles)
	b.TextStyles.ApplyStyles(defText, b.Styles)
	b.DrawingStyles.ApplyStyles(defDraw, b.Styles)
}

type TextStyles struct {
	lineHeight LineHeight
	width      Width
	hAlign     HAlign
}

func (ts *TextStyles) ApplyStyles(def TextStyles, styles Styles) {
	*ts = def
	for _, s := range styles {
		switch s := s.(type) {
		case LineHeight:
			ts.lineHeight = s
		case Width:
			ts.width = s
		case HAlign:
			ts.hAlign = s
		}
	}
}

type Text struct {
	Styled
	TextStyles
	XMLName xml.Name `xml:"Text"`
	Text    string   `xml:",chardata"`
}

func (t *Text) ApplyStyles(def TextStyles) {
	t.TextStyles = def
	for _, s := range t.Styles {
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
