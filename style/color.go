package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type RGB struct {
	R, G, B uint8
}

func (c RGB) String() string {
	return fmt.Sprintf("RGB(%d,%d,%d)", c.R, c.G, c.B)
}

func makeRGB(r, g, b uint8) RGB {
	return RGB{
		R: r,
		G: g,
		B: b,
	}
}

var Black RGB = RGB{0, 0, 0}
var White RGB = RGB{255, 255, 255}

func (c *RGB) UnmarshalStyle(s string) error {
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	if len(s) != 6 {
		return errors.Errorf("invalid color hex-string (%s)", s)
	}
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return errors.Wrapf(err, "parse color hex string (%s)", s)
	}
	c.R = uint8(n >> 16)
	c.G = uint8(n >> 8)
	c.B = uint8(n)
	return nil
}

type Color struct {
	Foreground RGB `style:"color"`
	Text       RGB `style:"text-color"`
	Background RGB `style:"background-color"`
}
