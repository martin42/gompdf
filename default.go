package gompdf

var DefaultFont = Font{
	fontFamily:     FontFamily("Arial"),
	fontPointSize:  FontPointSize(12),
	fontStyle:      FontStyleNormal,
	fontWeight:     FontWeightNormal,
	fontDecoration: FontDecorationNormal,
}

var DefaultText = Text{
	lineHeight: 1.5,
	width:      -1,
	hAlign:     HAlignLeft,
}
