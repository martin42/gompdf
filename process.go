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
		currFont: DefaultFont,
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
	p.processFont(&DefaultFont)

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
			p.processFont(i)
		case *LineFeed:
			p.processLineFeed(i)
		case *Box:
			p.renderBox(i)
		case *Text:
			p.renderText(i)
		}
	}
}

func (p *Processor) processFont(fnt *Font) {
	fnt.ApplyStyles(p.currFont)
	p.pdf.SetFont(string(fnt.fontFamily), fpdfFontStyle(fnt.fontStyle, fnt.fontWeight, fnt.fontDecoration), float64(fnt.fontPointSize))
	p.currFont = *fnt
}

func (p *Processor) processLineFeed(lf *LineFeed) {
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lf.Lines
	p.pdf.Ln(height)
}

func (p *Processor) renderText(text *Text) {
	text.ApplyStyles(DefaultText)
	var width float64
	if text.width > 0 {
		width = float64(text.width)
	} else {
		pw, _ := p.pdf.GetPageSize()
		l, _, r, _ := p.pdf.GetMargins()
		width = pw - (l + r) - 3
		Logf("set auto width: %.1f", width)
	}
	p.write(text.Text, width, float64(text.lineHeight), text.hAlign)
}

func (p *Processor) renderBox(box *Box) {

}
