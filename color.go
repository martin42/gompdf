package gompdf

import (
	"image/color"
	"strconv"
	"strings"
)

func RGBAFromHexColor(s string) color.RGBA {
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	if len(s) != 6 {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 0,
		}
	}
	rh := s[0:2]
	gh := s[2:4]
	bh := s[4:6]

	ru, _ := strconv.ParseUint(rh, 16, 8)
	gu, _ := strconv.ParseUint(gh, 16, 8)
	bu, _ := strconv.ParseUint(bh, 16, 8)

	cr := color.RGBA{
		R: uint8(ru),
		G: uint8(gu),
		B: uint8(bu),
	}
	return cr
}
