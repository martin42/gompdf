package gompdf

import (
	"io"

	"github.com/mazzegi/gompdf/style"
	"github.com/signintech/gopdf"
)

type GoPdfEngine struct {
	pdf *gopdf.GoPdf
}

func NewGoPdfEngine() *GoPdfEngine {
	e := &GoPdfEngine{
		pdf: &gopdf.GoPdf{},
	}
	return e
}

func (e *GoPdfEngine) Write(w io.Writer) error {
	return e.pdf.Write(w)
}

func (e *GoPdfEngine) Setup(c Config) {
	pc := gopdf.Config{
		Unit: gopdf.UnitMM,
	}
	switch c.Format {
	case FormatA3:
		pc.PageSize = *gopdf.PageSizeA3
	case FormatA4:
		pc.PageSize = *gopdf.PageSizeA4
	case FormatA5:
		pc.PageSize = *gopdf.PageSizeA5
	case FormatLegal:
		pc.PageSize = *gopdf.PageSizeLegal
	case FormatLetter:
		pc.PageSize = *gopdf.PageSizeLetter
	}
	e.pdf.Start(pc)
}

func (e *GoPdfEngine) AddPage() {
	e.pdf.AddPage()
}

func (e *GoPdfEngine) SetX(v float64) {
	e.pdf.SetX(v)
}

func (e *GoPdfEngine) SetY(v float64) {
	e.pdf.SetY(v)
}

func (e *GoPdfEngine) GetX() float64 {
	return e.pdf.GetX()
}

func (e *GoPdfEngine) GetY() float64 {
	return e.pdf.GetY()
}

func (e *GoPdfEngine) Margins() PageMargins {
	var pm PageMargins
	pm.Left, pm.Top, pm.Right, pm.Bottom = e.pdf.Margins()
	return pm
}

func (e *GoPdfEngine) UseFont(fnt style.Font) {
	e.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), int(fnt.PointSize))
}

func (e *GoPdfEngine) Image(src string, x, y, width, height float64) error {
	return e.pdf.Image(src, x, y, &gopdf.Rect{W: width, H: height})
}

func (e *GoPdfEngine) TextWidth(s string) float64 {
	w, err := e.pdf.MeasureTextWidth(s)
	if err != nil {
		return 0
	}
	return w
}

func (e *GoPdfEngine) SetTextColor(cr style.RGB) {
	e.pdf.SetTextColor(cr.R, cr.G, cr.B)
}

func (e *GoPdfEngine) Text(s string) {
	e.pdf.Text(s)
}

//Drawing stuff
func (e *GoPdfEngine) Line(x0, y0, x1, y1 float64)          {}
func (e *GoPdfEngine) SetLineWidth(float64)                 {}
func (e *GoPdfEngine) SetStrokeColor(style.RGB)             {}
func (e *GoPdfEngine) SetFillColor(style.RGB)               {}
func (e *GoPdfEngine) FillRect(x, y, width, height float64) {} //FillFromUpperLeftWithStyle "F"
