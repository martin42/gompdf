package gompdf

import (
	"fmt"
	"io"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/martin42/gompdf/style"
)

type ProcessOption func(p *Processor) error

type Processor struct {
	doc *Document
	pdf *gofpdf.Fpdf

	fontDir  string
	codePage string

	translateUnicode func(string) string

	currStyles style.Styles
}

func WithFontDir(dir string) ProcessOption {
	return func(p *Processor) error {
		p.fontDir = dir
		return nil
	}
}

func WithCodePage(cp string) ProcessOption {
	return func(p *Processor) error {
		p.codePage = cp
		return nil
	}
}

func NewProcessor(doc *Document, options ...ProcessOption) (*Processor, error) {
	p := &Processor{
		doc:        doc,
		fontDir:    "fonts",
		codePage:   "",
		currStyles: DefaultStyle,
	}
	for _, o := range options {
		err := o(p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Processor) Process(w io.Writer) error {
	start := time.Now()
	fmt.Printf("run instructions ...\n")
	p.pdf = gofpdf.New(
		fpdfOrientation(p.doc.Default.Orientation),
		fpdfUnit(p.doc.Default.Unit),
		fpdfFormat(p.doc.Default.Format),
		p.fontDir,
	)
	p.translateUnicode = p.pdf.UnicodeTranslatorFromDescriptor(p.codePage)
	p.pdf.SetHeaderFunc(func() {
		p.processInstructions(p.doc.Header)
	})
	p.pdf.SetFooterFunc(func() {
		p.processInstructions(p.doc.Footer)
	})
	p.applyDefaults()
	p.applyFont(p.currStyles.Font)

	p.pdf.AddPage()
	p.processInstructions(p.doc.Body)

	err := p.pdf.Error()
	if err != nil {
		return err
	}
	fmt.Printf("run instructions ... in (%s)\n", time.Since(start))
	return p.pdf.Output(w)
}

func (p *Processor) applyDefaults() {
	p.pdf.SetAutoPageBreak(p.doc.Default.PageBreaks == PageBreakModeAuto, p.doc.Default.PageMargins.Bottom)
	p.pdf.SetMargins(p.doc.Default.PageMargins.Left, p.doc.Default.PageMargins.Top, p.doc.Default.PageMargins.Right)
}

func (p *Processor) appliedStyles(i Instruction) style.Styles {
	st := p.currStyles
	i.Apply(p.doc.styleClasses, &st)
	return st
}

func (p *Processor) processInstructions(is Instructions) {
	for _, i := range is.iss {
		//i.Apply(p.doc.styleClasses, &styles)
		switch i := i.(type) {
		case *Font:
			i.Apply(p.doc.styleClasses, &p.currStyles)
			p.applyFont(p.currStyles.Font)
		case *LineFeed:
			p.processLineFeed(i, p.appliedStyles(i))
		case *SetX:
			p.pdf.SetX(i.X)
		case *SetY:
			p.pdf.SetY(i.Y)
		case *SetXY:
			p.pdf.SetXY(i.X, i.Y)
		case *Box:
			p.renderBox(i, p.appliedStyles(i))
		case *Text:
			p.renderText(i, p.appliedStyles(i))
		case *Table:
			p.renderTable(i, p.appliedStyles(i))
		}
	}
}

func (p *Processor) effectiveWidth(width float64) float64 {
	if width > 0 {
		return width
	}
	pw, _ := p.pdf.GetPageSize()
	l, _, r, _ := p.pdf.GetMargins()
	return pw - (l + r) - 3 // without substracting 3 it doesn't fit
}

func (p *Processor) applyFont(fnt style.Font) {
	p.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), float64(fnt.PointSize))
}

func (p *Processor) processLineFeed(lf *LineFeed, sty style.Styles) {
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lf.Lines
	p.pdf.Ln(height)
}

func (p *Processor) renderText(text *Text, sty style.Styles) {
	p.write(text.Text, p.effectiveWidth(sty.Dimension.Width), sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font)
}

func (p *Processor) renderBox(box *Box, sty style.Styles) {

	x0, y0 := p.pdf.GetXY()
	width := p.effectiveWidth(sty.Dimension.Width)
	textWidth := width - sty.Box.Padding.Left - sty.Box.Padding.Right - 2 //without -2 it writes over the border

	height := p.textHeight(box.Text, textWidth, sty.Dimension.LineHeight, sty.Font)
	y1 := y0 + height + sty.Box.Padding.Bottom

	p.pdf.SetLineWidth(sty.Draw.LineWidth)
	p.pdf.SetDrawColor(int(sty.Color.Foreground.R), int(sty.Color.Foreground.G), int(sty.Color.Foreground.B))
	p.pdf.SetFillColor(int(sty.Color.Background.R), int(sty.Color.Background.G), int(sty.Color.Background.B))

	p.pdf.Rect(x0, y0, width, y1-y0, "F")
	p.pdf.MoveTo(x0, y0)
	if sty.Box.Border.Top > 0 {
		p.pdf.LineTo(x0+width, y0)
	} else {
		p.pdf.MoveTo(x0+width, y0)
	}
	if sty.Box.Border.Right > 0 {
		p.pdf.LineTo(x0+width, y1)
	} else {
		p.pdf.MoveTo(x0+width, y1)
	}
	if sty.Box.Border.Bottom > 0 {
		p.pdf.LineTo(x0, y1)
	} else {
		p.pdf.MoveTo(x0, y1)
	}
	if sty.Box.Border.Left > 0 {
		p.pdf.LineTo(x0, y0)
	} else {
		p.pdf.MoveTo(x0, y0)
	}
	p.pdf.DrawPath("D")
	p.pdf.SetXY(x0, y1)

	p.pdf.SetY(y0 + sty.Box.Padding.Top)
	p.pdf.SetX(x0 + sty.Box.Padding.Left)
	p.write(box.Text, textWidth, sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font)
}
