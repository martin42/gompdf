package gompdf

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

func fpdfFontStyle(fs FontStyle, fw FontWeight, fd FontDecoration) string {
	s := ""
	switch fs {
	case FontStyleItalic:
		s += "I"
	}
	switch fw {
	case FontWeightBold:
		s += "B"
	}
	switch fd {
	case FontDecorationUnderline:
		s += "U"
	}
	return s
}
