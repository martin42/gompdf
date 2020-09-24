package gompdf

import (
	"io"

	"github.com/mazzegi/gompdf/style"
	"github.com/signintech/gopdf"
)

type GoPdfTTFFont struct {
	Name    string
	TTFPath string
	Option  gopdf.TtfOption
}

type GoPdfEngine struct {
	pdf   *gopdf.GoPdf
	fonts []GoPdfTTFFont
}

func NewGoPdfEngine(fonts ...GoPdfTTFFont) *GoPdfEngine {
	e := &GoPdfEngine{
		pdf:   &gopdf.GoPdf{},
		fonts: fonts,
	}
	return e
}

func (e *GoPdfEngine) Write(w io.Writer) error {
	Logf("engine: write")
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
	Logf("engine: start (cfg.page-size=%v)", pc.PageSize)
	e.pdf.Start(pc)

	//fonts
	for _, f := range e.fonts {
		e.pdf.AddTTFFontWithOption(f.Name, f.TTFPath, f.Option)
	}

	Logf("engine: set margins %v", c.Margins)
	e.pdf.SetMargins(c.Margins.Left, c.Margins.Top, c.Margins.Right, c.Margins.Bottom)
}

func (e *GoPdfEngine) AddPage() {
	Logf("engine: add-page")
	e.pdf.AddPage()
}

func (e *GoPdfEngine) SetX(v float64) {
	Logf("engine: set-x (%f)", v)
	e.pdf.SetX(v)
}

func (e *GoPdfEngine) SetY(v float64) {
	Logf("engine: set-y (%f)", v)
	e.pdf.SetY(v)
}

func (e *GoPdfEngine) GetX() float64 {
	Logf("engine: get-x")
	return e.pdf.GetX()
}

func (e *GoPdfEngine) GetY() float64 {
	Logf("engine: get-y")
	return e.pdf.GetY()
}

func (e *GoPdfEngine) Margins() PageMargins {
	Logf("engine: margins")
	var pm PageMargins
	pm.Left, pm.Top, pm.Right, pm.Bottom = e.pdf.Margins()
	return pm
}

func (e *GoPdfEngine) UseFont(fnt style.Font) {
	Logf("engine: use-font (%q, %q, %d)", fnt.Family, fpdfFontStyle(fnt), int(fnt.PointSize))
	e.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), int(fnt.PointSize))
}

func (e *GoPdfEngine) Image(src string, x, y, width, height float64) error {
	Logf("engine: image tbd")
	return e.pdf.Image(src, x, y, &gopdf.Rect{W: width, H: height})
}

func (e *GoPdfEngine) TextWidth(s string) float64 {
	Logf("engine: measure-text.width of %q", s)
	w, err := e.pdf.MeasureTextWidth(s)
	if err != nil {
		return 0
	}
	return w
}

func (e *GoPdfEngine) SetTextColor(cr style.RGB) {
	Logf("engine: set-text-color %q", cr.String())
	e.pdf.SetTextColor(cr.R, cr.G, cr.B)
}

func (e *GoPdfEngine) Text(s string) {
	Logf("engine: text %q", s)
	e.pdf.Text(s)
}

//Drawing stuff
func (e *GoPdfEngine) Line(x0, y0, x1, y1 float64) {
	e.pdf.Line(x0, y0, x1, y1)
}

func (e *GoPdfEngine) SetLineWidth(v float64) {
	e.pdf.SetLineWidth(v)
}

func (e *GoPdfEngine) SetStrokeColor(cr style.RGB) {
	e.pdf.SetStrokeColor(cr.R, cr.G, cr.B)
}

func (e *GoPdfEngine) SetFillColor(cr style.RGB) {
	e.pdf.SetFillColor(cr.R, cr.G, cr.B)
}

func (e *GoPdfEngine) FillRect(x, y, width, height float64) {
	e.pdf.RectFromUpperLeftWithStyle(x, y, width, height, "F")
}
