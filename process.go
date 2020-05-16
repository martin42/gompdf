package gompdf

import (
	"fmt"
	"io"
	"time"

	//"github.com/jung-kurt/gofpdf/v2"
	"github.com/mazzegi/gompdf/style"
	gofpdf "github.com/signintech/gopdf"
)

type ProcessOption func(p *Processor) error

type Processor struct {
	doc *Document
	pdf *gofpdf.GoPdf

	fontDir  string
	codePage string
	pageSize *gofpdf.Rect

	transformText func(string) string

	currStyles style.Styles
	currFont   style.Font
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
	p.pageSize = gofpdf.PageSizeA4
	p.pdf = &gofpdf.GoPdf{}
	p.pdf.Start(gofpdf.Config{
		PageSize: *p.pageSize,
	})
	p.pdf.AddPage()
	err := p.pdf.AddTTFFont("times", "fonts/times.ttf")
	if err != nil {
		return err
	}

	//TODO: Units
	//TODO: Orientation
	//TODO: Header/Footer
	//TODO: Page Numbering
	//TODO: Auto Page break

	// gofpdf.New(
	// 	fpdfOrientation(p.doc.Default.Orientation),
	// 	fpdfUnit(p.doc.Default.Unit),
	// 	fpdfFormat(p.doc.Default.Format),
	// 	p.fontDir,
	// )

	// p.pdf.AliasNbPages("{np}")
	// translateUnicode := p.pdf.UnicodeTranslatorFromDescriptor(p.codePage)
	p.transformText = func(s string) string {
		// ts := strings.Replace(s, "{cp}", fmt.Sprintf("%d", p.pdf.PageNo()), -1)
		// return translateUnicode(ts)
		return s
	}

	// p.pdf.SetHeaderFunc(func() {
	// 	p.processInstructions(p.doc.Header)
	// })
	// p.pdf.SetFooterFunc(func() {
	// 	p.processInstructions(p.doc.Footer)
	// })
	p.applyDefaults()
	p.applyFont(p.currStyles.Font)

	p.pdf.AddPage()
	p.processInstructions(p.doc.Body)

	// err := p.pdf.Error()
	// if err != nil {
	// 	return err
	// }
	fmt.Printf("run instructions ... in (%s)\n", time.Since(start))
	return p.pdf.Write(w)
	//return p.pdf.Output(w)
}

func (p *Processor) applyDefaults() {
	//p.pdf.SetAutoPageBreak(p.doc.Default.PageBreaks == PageBreakModeAuto, p.doc.Default.PageMargins.Bottom)
	p.pdf.SetMargins(p.doc.Default.PageMargins.Left, p.doc.Default.PageMargins.Top, p.doc.Default.PageMargins.Right, p.doc.Default.PageMargins.Bottom)
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
			p.appliedStyles(i)
			p.processLineFeed(i)
		case *SetX:
			p.pdf.SetX(i.X)
		case *SetY:
			p.pdf.SetY(i.Y)
		case *SetXY:
			p.pdf.SetX(i.X)
			p.pdf.SetY(i.Y)
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
	pw := p.pageSize.W
	//pw, _ := p.pdf.GetPageSize()
	l, _, r, _ := p.pdf.Margins()
	return pw - (l + r) - 3 // without substracting 3 it doesn't fit
}

func (p *Processor) applyFont(fnt style.Font) {
	p.currFont = fnt
	p.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), int(fnt.PointSize))
}

func (p *Processor) processLineFeed(lf *LineFeed) {
	//_, fontHeight := p.pdf.GetFontSize()
	fontHeight := p.currFont.PointSize
	height := fontHeight * lf.Lines
	p.pdf.SetY(p.pdf.GetY() + height)
}

func (p *Processor) ln(h float64) {
	p.pdf.SetY(p.pdf.GetY() + h)
	p.pdf.SetX(p.pdf.MarginLeft())
}

func (p *Processor) renderText(text *Text, sty style.Styles) {
	p.write(text.Text, p.effectiveWidth(sty.Dimension.Width), sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font, sty.Color.Text)
}

func (p *Processor) GetXY() (float64, float64) {
	return p.pdf.GetX() + p.pdf.MarginLeft(), p.pdf.GetY()
}

func (p *Processor) SetXY(x, y float64) {
	p.pdf.SetX(x)
	p.pdf.SetY(y)
}

func (p *Processor) renderTextBox(text string, sty style.Styles) {
	width := p.effectiveWidth(sty.Dimension.Width)
	textWidth := width - sty.Box.Padding.Left - sty.Box.Padding.Right - 3 //without -2 it writes over the border

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

	x0, y0 := p.GetXY()
	Logf("render-box: x,y: %f, %f", x0, y0)
	//_, ph := p.pdf.GetPageSize()
	ph := p.pageSize.H

	if y0+height >= ph {
		p.pdf.AddPage()
		x0, y0 = p.GetXY()
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
	//p.pdf.Ln(sty.Dimension.LineHeight + sty.Box.Padding.Bottom)
	p.ln(sty.Dimension.LineHeight + sty.Box.Padding.Bottom)
}

func (p *Processor) renderImage(img *Image, sty style.Styles) {
	x0, y0 := p.GetXY()
	x0 += sty.Dimension.OffsetX
	y0 += sty.Dimension.OffsetY
	//p.pdf.ImageOptions(img.Source, x0, y0, sty.Dimension.Width, sty.Dimension.Height, false, gofpdf.ImageOptions{}, 0, "")
	r := gofpdf.Rect{
		W: sty.Dimension.Width,
		H: sty.Dimension.Height,
	}
	p.pdf.Image(img.Source, x0, y0, &r)
}
