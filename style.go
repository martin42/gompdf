package gompdf

import (
	"bytes"
	"fmt"
	"image/color"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type Style interface{}

type Styles []Style

func (s Styles) FindByPrototype(proto Style) (Style, bool) {
	for _, st := range s {
		if reflect.TypeOf(proto) == reflect.TypeOf(st) {
			return st, true
		}
	}
	return nil, false
}

func (s Styles) Get(proto *Style) bool {
	for _, st := range s {
		if reflect.TypeOf(*proto) == reflect.TypeOf(st) {
			*proto = st
			return true
		}
	}
	return false
}

type StyleClass struct {
	Name   string
	Styles Styles
}

type StyleClasses []StyleClass

func (scs StyleClasses) Find(name string) (StyleClass, bool) {
	for _, sc := range scs {
		if sc.Name == name {
			return sc, true
		}
	}
	return StyleClass{}, false
}

type MakeStyleFnc func(value string) (Style, error)

var styleRegistry *StyleRegistry

func init() {
	styleRegistry = NewStyleRegistry()
	styleRegistry.Register("border", makeBorder)
	styleRegistry.Register("padding", makePadding)
	styleRegistry.Register("margin", makeMargin)
	styleRegistry.Register("position", makePosition)
	styleRegistry.Register("width", makeWidth)
	styleRegistry.Register("height", makeHeight)
	styleRegistry.Register("line-height", makeLineHeight)
	styleRegistry.Register("font-family", makeFontFamily)
	styleRegistry.Register("font-size", makeFontPointSize)
	styleRegistry.Register("font-style", makeFontStyle)
	styleRegistry.Register("font-weight", makeFontWeight)
	styleRegistry.Register("font-decoration", makeFontDecoration)
	styleRegistry.Register("h-align", makeHAlign)
	styleRegistry.Register("v-align", makeVAlign)
	styleRegistry.Register("background-color", makeBackgroundColor)
	styleRegistry.Register("color", makeColor)
	styleRegistry.Register("line-width", makeLineWidth)
}

type StyleRegistry struct {
	types map[string]MakeStyleFnc
}

func NewStyleRegistry() *StyleRegistry {
	return &StyleRegistry{
		types: map[string]MakeStyleFnc{},
	}
}

func (r *StyleRegistry) Register(name string, makeFnc MakeStyleFnc) {
	r.types[name] = makeFnc
}

func (r *StyleRegistry) Decode(name string, value string) (Style, error) {
	makeFnc, contains := r.types[name]
	if !contains {
		return nil, errors.Errorf("registry-decode: (%s) is not registered", name)
	}
	sty, err := makeFnc(value)
	if err != nil {
		return nil, err
	}
	return sty, nil
}

const whitespace = " \r\n\t"

func ParseClasses(bs []byte) (StyleClasses, error) {
	bs = bytes.Replace(bs, []byte("\r"), []byte(" "), -1)
	bs = bytes.Replace(bs, []byte("\n"), []byte(" "), -1)
	bs = bytes.Replace(bs, []byte("\t"), []byte(" "), -1)
	scs := StyleClasses{}
	pos := 0
	for {
		curr := bs[pos:]
		i := bytes.IndexByte(curr, '{')
		if i < 0 {
			return scs, nil
		}
		name := bytes.Trim(curr[:i], whitespace)
		if len(name) == 0 {
			return nil, errors.Errorf("style class without name")
		}
		in := bytes.IndexByte(curr[i:], '}')
		if in < 0 {
			return nil, errors.Errorf("non matching brace")
		}
		in += i
		styles, err := ParseStyles(curr[i+1 : in])
		if err != nil {
			return nil, errors.Wrap(err, "parse style")
		}
		scs = append(scs, StyleClass{
			Name:   string(name),
			Styles: styles,
		})
		pos += in + 1
	}
}

func ParseStyles(s []byte) (Styles, error) {
	styles := Styles{}
	iss := bytes.Split(s, []byte(";"))
	for _, i := range iss {
		i = bytes.Trim(i, " \r\n\t")
		if len(i) == 0 {
			continue
		}
		is := bytes.Split(i, []byte(":"))
		if len(is) != 2 {
			return nil, errors.Errorf("invalid style syntax (%s)", i)
		}
		key := bytes.Trim(is[0], " \r\n\t")
		val := bytes.Trim(is[1], " \r\n\t")
		sty, err := styleRegistry.Decode(string(key), string(val))
		if err != nil {
			return nil, errors.Wrapf(err, "decode style (%s: %s)", key, val)
		}
		if _, contains := styles.FindByPrototype(sty); contains {
			return nil, errors.Errorf("style (%s: %s) is already set", key, val)
		}
		styles = append(styles, sty)
	}
	return styles, nil
}

type Border struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func makeBorder(value string) (Style, error) {
	b := Border{}
	_, err := fmt.Fscanf(bytes.NewBufferString(value), "%d,%d,%d,%d", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return nil, errors.Wrapf(err, "scan border value (%s)", value)
	}
	return b, nil
}

type Padding struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func makePadding(value string) (Style, error) {
	b := Padding{}
	_, err := fmt.Fscanf(bytes.NewBufferString(value), "%f,%f,%f,%f", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return nil, errors.Wrapf(err, "scan padding value (%s)", value)
	}
	return b, nil
}

type Margin struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func makeMargin(value string) (Style, error) {
	b := Margin{}
	_, err := fmt.Fscanf(bytes.NewBufferString(value), "%f,%f,%f,%f", &b.Left, &b.Top, &b.Right, &b.Bottom)
	if err != nil {
		return nil, errors.Wrapf(err, "scan margin value (%s)", value)
	}
	return b, nil
}

type Position struct {
	X float64
	Y float64
}

func makePosition(value string) (Style, error) {
	b := Position{}
	_, err := fmt.Fscanf(bytes.NewBufferString(value), "%f,%f", &b.X, &b.Y)
	if err != nil {
		return nil, errors.Wrapf(err, "scan position value (%s)", value)
	}
	return b, nil
}

type FontFamily string

func makeFontFamily(value string) (Style, error) {
	return FontFamily(value), nil
}

type FontPointSize float64

func makeFontPointSize(value string) (Style, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "parse float (%s)", value)
	}
	return FontPointSize(v), nil
}

type FontStyle string

const (
	FontStyleNormal FontStyle = "normal"
	FontStyleItalic FontStyle = "italic"
)

func makeFontStyle(value string) (Style, error) {
	switch value {
	case string(FontStyleNormal):
		return FontStyleNormal, nil
	case string(FontStyleItalic):
		return FontStyleItalic, nil
	}
	return nil, errors.Errorf("invalid font style (%s)", value)
}

type FontWeight string

const (
	FontWeightNormal FontWeight = "normal"
	FontWeightBold   FontWeight = "bold"
)

func makeFontWeight(value string) (Style, error) {
	switch value {
	case string(FontWeightNormal):
		return FontWeightNormal, nil
	case string(FontWeightBold):
		return FontWeightBold, nil
	}
	return nil, errors.Errorf("invalid font weight (%s)", value)
}

type FontDecoration string

const (
	FontDecorationNormal    FontDecoration = "normal"
	FontDecorationUnderline FontDecoration = "underline"
)

func makeFontDecoration(value string) (Style, error) {
	switch value {
	case string(FontDecorationNormal):
		return FontDecorationNormal, nil
	case string(FontDecorationUnderline):
		return FontDecorationUnderline, nil
	}
	return nil, errors.Errorf("invalid font decoration (%s)", value)
}

type HAlign string

const (
	HAlignLeft   HAlign = "left"
	HAlignRight  HAlign = "right"
	HAlignCenter HAlign = "center"
)

func makeHAlign(value string) (Style, error) {
	switch value {
	case string(HAlignLeft):
		return HAlignLeft, nil
	case string(HAlignRight):
		return HAlignRight, nil
	case string(HAlignCenter):
		return HAlignCenter, nil
	}
	return nil, errors.Errorf("invalid halign value (%s)", value)
}

type VAlign string

const (
	VAlignTop    VAlign = "top"
	VAlignMiddle VAlign = "middle"
	VAlignBottom VAlign = "bottom"
)

func makeVAlign(value string) (Style, error) {
	switch value {
	case string(VAlignTop):
		return VAlignTop, nil
	case string(VAlignMiddle):
		return VAlignMiddle, nil
	case string(VAlignBottom):
		return VAlignBottom, nil
	}
	return nil, errors.Errorf("invalid valign value (%s)", value)
}

type Width float64

func makeWidth(value string) (Style, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "parse width (%s)", value)
	}
	return Width(v), nil
}

type Height float64

func makeHeight(value string) (Style, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "parse height (%s)", value)
	}
	return Height(v), nil
}

type ImageFlow int

const (
	ImageFlowEnabled  ImageFlow = 1
	ImageFlowDisabled ImageFlow = 0
)

func makeImageFlow(value string) (Style, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parse image flow (%s)", value)
	}
	return ImageFlow(v), nil
}

type LineHeight float64

func makeLineHeight(value string) (Style, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "parse line-height (%s)", value)
	}
	return LineHeight(v), nil
}

type BackgroundColor color.RGBA

func makeBackgroundColor(value string) (Style, error) {
	return BackgroundColor(RGBAFromHexColor(value)), nil
}

type Color color.RGBA

func makeColor(value string) (Style, error) {
	return Color(RGBAFromHexColor(value)), nil
}

type LineWidth float64

func makeLineWidth(value string) (Style, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "parse line-width (%s)", value)
	}
	return LineWidth(v), nil
}

//type
// crOut := RGBAFromHexColor(r.Outline)
// 	crFill := RGBAFromHexColor(r.Fill)
// 	p.pdf.SetDrawColor(int(crOut.R), int(crOut.G), int(crOut.B))
// 	p.pdf.SetFillColor(int(crFill.R), int(crFill.G), int(crFill.B))
