package gompdf

import "github.com/martin42/gompdf/style"

func fpdfOrientation(o Orientation) string {
	switch o {
	case OrientationPortrait:
		return "P"
	case OrientationLandscape:
		return "L"
	}
	return ""
}

func fpdfUnit(u Unit) string {
	return string(u)
}

func fpdfFormat(f Format) string {
	return string(f)
}

func fpdfFontStyle(fnt style.Font) string {
	s := ""
	switch fnt.Style {
	case style.FontStyleItalic:
		s += "I"
	}
	switch fnt.Weight {
	case style.FontWeightBold:
		s += "B"
	}
	switch fnt.Decoration {
	case style.FontDecorationUnderline:
		s += "U"
	}
	return s
}
