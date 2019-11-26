package gompdf

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mazzegi/gompdf/style"
)

type ProcessOption func(p *Processor) error

type Processor struct {
	doc *Document
	pdf *gofpdf.Fpdf

	fontDir  string
	codePage string

	transformText func(string) string

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
	p.pdf.AliasNbPages("{np}")
	translateUnicode := p.pdf.UnicodeTranslatorFromDescriptor(p.codePage)
	p.transformText = func(s string) string {
		ts := strings.Replace(s, "{cp}", fmt.Sprintf("%d", p.pdf.PageNo()), -1)
		return translateUnicode(ts)
	}

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
			p.renderTextBox(i.Text, p.appliedStyles(i))
		case *Text:
			p.renderText(i, p.appliedStyles(i))
		case *Table:
			p.renderTable(i, p.appliedStyles(i))
		case *Image:
			p.renderImage(i, p.appliedStyles(i))
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
	//Logf("set-font: %s, %s, %.1f", fnt.Family, fpdfFontStyle(fnt), fnt.PointSize)
	p.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), float64(fnt.PointSize))
}

func (p *Processor) processLineFeed(lf *LineFeed, sty style.Styles) {
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lf.Lines
	p.pdf.Ln(height)
}

func (p *Processor) renderText(text *Text, sty style.Styles) {
	p.write(text.Text, p.effectiveWidth(sty.Dimension.Width), sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font, sty.Color.Text)
}

func (p *Processor) renderTextBox(text string, sty style.Styles) {
	width := p.effectiveWidth(sty.Dimension.Width)
	textWidth := width - sty.Box.Padding.Left - sty.Box.Padding.Right - 2 //without -2 it writes over the border

	var height float64
	if sty.Dimension.Height < 0 {
		if text != "" {
			height = p.textHeight(text, textWidth, sty.Dimension.LineHeight, sty.Font)
		} else {
			height = p.textHeight("Üg", textWidth, sty.Dimension.LineHeight, sty.Font)
		}
	} else {
		height = sty.Dimension.Height
	}

	x0, y0 := p.pdf.GetXY()
	_, ph := p.pdf.GetPageSize()
	if y0+height >= ph {
		p.pdf.AddPage()
		x0, y0 = p.pdf.GetXY()
	}

	x0 += sty.Dimension.OffsetX
	y0 += sty.Dimension.OffsetY
	y1 := y0 + height + sty.Box.Padding.Top + sty.Box.Padding.Bottom
	x1 := x0 + width
	p.drawBox(x0, y0, x1, y1, sty)

	//Reset, to start writing at top left
	p.pdf.SetY(y0 + sty.Box.Padding.Top)
	p.pdf.SetX(x0 + sty.Box.Padding.Left)
	p.write(text, textWidth, sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font, sty.Color.Text)
	p.pdf.Ln(sty.Dimension.LineHeight + sty.Box.Padding.Bottom)
}

func (p *Processor) renderImage(img *Image, sty style.Styles) {
	x0, y0 := p.pdf.GetXY()
	x0 += sty.Dimension.OffsetX
	y0 += sty.Dimension.OffsetY
	p.pdf.ImageOptions(img.Source, x0, y0, sty.Dimension.Width, sty.Dimension.Height, false, gofpdf.ImageOptions{}, 0, "")
}
