package style

import "fmt"

type FontStyle string

const (
	FontStyleNormal FontStyle = "normal"
	FontStyleItalic FontStyle = "italic"
)

type FontWeight string

const (
	FontWeightNormal FontWeight = "normal"
	FontWeightBold   FontWeight = "bold"
)

type FontDecoration string

const (
	FontDecorationNormal    FontDecoration = "normal"
	FontDecorationUnderline FontDecoration = "underline"
)

type Font struct {
	Family     string         `style:"font-family"`
	PointSize  float64        `style:"font-point-size"`
	Style      FontStyle      `style:"font-style"`
	Weight     FontWeight     `style:"font-weight"`
	Decoration FontDecoration `style:"font-decoration"`
}

func (f Font) String() string {
	return fmt.Sprintf("%s(%f): style=%q weight=%q dec=%q", f.Family, f.PointSize, f.Style, f.Weight, f.Decoration)
}
