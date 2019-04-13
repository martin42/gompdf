package style

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
