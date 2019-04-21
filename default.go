package gompdf

import (
	"github.com/martin42/gompdf/style"
)

var DefaultStyle = style.Styles{
	Font: style.Font{
		Family:     "Arial",
		PointSize:  12,
		Style:      style.FontStyleNormal,
		Weight:     style.FontWeightNormal,
		Decoration: style.FontDecorationNormal,
	},
	Box: style.Box{
		Border:  style.Border{Left: 0, Top: 0, Right: 0, Bottom: 0},
		Padding: style.Padding{Left: 0, Top: 0, Right: 0, Bottom: 0},
		Margin:  style.Margin{Left: 0, Top: 0, Right: 0, Bottom: 0},
	},
	Dimension: style.Dimension{
		Width:       -1,
		Height:      -1,
		ColumnWidth: -1,
		LineHeight:  1.5,
	},
	Align: style.Align{
		HAlign: style.HAlignLeft,
		VAlign: style.VAlignTop,
	},
	Color: style.Color{
		Foreground: style.Black,
		Background: style.White,
	},
}

// var DefaultFontStyles = FontStyles{
// 	fontFamily:     FontFamily("Arial"),
// 	fontPointSize:  FontPointSize(12),
// 	fontStyle:      FontStyleNormal,
// 	fontWeight:     FontWeightNormal,
// 	fontDecoration: FontDecorationNormal,
// }

// var DefaultTextStyles = TextStyles{
// 	lineHeight: 1.5,
// 	width:      -1,
// 	hAlign:     HAlignLeft,
// }

// var DefaultBoxStyles = BoxStyles{
// 	border: Border{
// 		Left:   1,
// 		Top:    1,
// 		Right:  1,
// 		Bottom: 1,
// 	},
// 	padding: Padding{
// 		Left:   2,
// 		Top:    2,
// 		Right:  2,
// 		Bottom: 2,
// 	},
// }

// var DefaultDrawingStyles = DrawingStyles{
// 	backgroundColor: BackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 0}),
// 	color:           Color(color.RGBA{R: 0, G: 0, B: 0, A: 0}),
// 	lineWidth:       0.2,
// }
