package gompdf

import (
	"fmt"
	"io"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ProcessOption func(p *Processor) error

type Processor struct {
	doc *Document
	pdf *gofpdf.Fpdf

	fontDir  string
	codePage string

	translateUnicode func(string) string

	currFont Font
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
		doc:      doc,
		fontDir:  "fonts",
		codePage: "",
		currFont: Font{
			FontStyles: DefaultFontStyles,
		},
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
	p.applyFont(&p.currFont)

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

func (p *Processor) processInstructions(is Instructions) {
	for _, i := range is {
		switch i := i.(type) {
		case *Font:
			p.applyFont(i)
		case *LineFeed:
			p.processLineFeed(i)
		case *SetX:
			p.pdf.SetX(i.X)
		case *SetY:
			p.pdf.SetY(i.Y)
		case *SetXY:
			p.pdf.SetXY(i.X, i.Y)
		case *Box:
			p.renderBox(i)
		case *Text:
			p.renderText(i)
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

func (p *Processor) applyFont(fnt *Font) {
	fnt.ApplyStyles(p.currFont.FontStyles)
	p.pdf.SetFont(string(fnt.fontFamily), fpdfFontStyle(fnt.fontStyle, fnt.fontWeight, fnt.fontDecoration), float64(fnt.fontPointSize))
	p.currFont = *fnt
}

func (p *Processor) processLineFeed(lf *LineFeed) {
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lf.Lines
	p.pdf.Ln(height)
}

func (p *Processor) renderText(text *Text) {
	text.ApplyStyles(DefaultTextStyles)
	p.write(text.Text, p.effectiveWidth(float64(text.width)), float64(text.lineHeight), text.hAlign)
}

func (p *Processor) renderBox(box *Box) {
	box.ApplyStyles(DefaultBoxStyles, DefaultTextStyles, DefaultDrawingStyles)
	x0, y0 := p.pdf.GetXY()
	width := p.effectiveWidth(float64(box.width))
	textWidth := width - box.padding.Left - box.padding.Right - 2 //without -2 it writes over the border

	height := p.textHeight(box.Text, textWidth, float64(box.lineHeight))
	y1 := y0 + height + box.padding.Bottom

	p.pdf.SetLineWidth(float64(box.lineWidth))
	p.pdf.SetDrawColor(int(box.color.R), int(box.color.G), int(box.color.B))
	p.pdf.SetFillColor(int(box.backgroundColor.R), int(box.backgroundColor.G), int(box.backgroundColor.B))

	p.pdf.Rect(x0, y0, width, y1-y0, "F")
	p.pdf.MoveTo(x0, y0)
	if box.border.Top > 0 {
		p.pdf.LineTo(x0+width, y0)
	} else {
		p.pdf.MoveTo(x0+width, y0)
	}
	if box.border.Right > 0 {
		p.pdf.LineTo(x0+width, y1)
	} else {
		p.pdf.MoveTo(x0+width, y1)
	}
	if box.border.Bottom > 0 {
		p.pdf.LineTo(x0, y1)
	} else {
		p.pdf.MoveTo(x0, y1)
	}
	if box.border.Left > 0 {
		p.pdf.LineTo(x0, y0)
	} else {
		p.pdf.MoveTo(x0, y0)
	}
	p.pdf.DrawPath("D")
	p.pdf.SetXY(x0, y1)

	p.pdf.SetY(y0 + box.padding.Top)
	p.pdf.SetX(x0 + box.padding.Left)
	p.write(box.Text, textWidth, float64(box.lineHeight), box.hAlign)
}
