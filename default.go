package gompdf

import "image/color"

var DefaultFontStyles = FontStyles{
	fontFamily:     FontFamily("Arial"),
	fontPointSize:  FontPointSize(12),
	fontStyle:      FontStyleNormal,
	fontWeight:     FontWeightNormal,
	fontDecoration: FontDecorationNormal,
}

var DefaultTextStyles = TextStyles{
	lineHeight: 1.5,
	width:      -1,
	hAlign:     HAlignLeft,
}

var DefaultBoxStyles = BoxStyles{
	border: Border{
		Left:   1,
		Top:    1,
		Right:  1,
		Bottom: 1,
	},
	padding: Padding{
		Left:   2,
		Top:    2,
		Right:  2,
		Bottom: 2,
	},
}

var DefaultDrawingStyles = DrawingStyles{
	backgroundColor: BackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 0}),
	color:           Color(color.RGBA{R: 0, G: 0, B: 0, A: 0}),
	lineWidth:       0.2,
}
