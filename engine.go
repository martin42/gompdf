package gompdf

import (
	"io"

	"github.com/mazzegi/gompdf/style"
)

type Config struct {
	Format  Format
	Margins PageMargins
}

type Engine interface {
	Setup(c Config) //Start
	AddPage()
	SetX(float64)
	SetY(float64)
	GetX() float64
	GetY() float64
	Margins() PageMargins

	UseFont(style.Font)
	Image(src string, x, y, width, height float64) error
	TextWidth(string) float64
	SetTextColor(style.RGB)
	Text(string)

	Line(x0, y0, x1, y1 float64)
	SetLineWidth(float64)
	SetStrokeColor(style.RGB)
	SetFillColor(style.RGB)
	FillRect(x, y, width, height float64) //FillFromUpperLeftWithStyle "F"

	Write(w io.Writer) error
}
