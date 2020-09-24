package gompdf

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	//"github.com/jung-kurt/gofpdf/v2"
	"github.com/mazzegi/gompdf/style"
	gofpdf "github.com/signintech/gopdf"
)

type ProcessOption func(p *Processor) error

type Processor struct {
	doc *Document
	//pdf *gofpdf.GoPdf
	pdf Engine

	fontDir  string
	codePage string
	pageSize *gofpdf.Rect

	transformText func(string) string

	currStyles       style.Styles
	currFont         style.Font
	currPage         int
	preventPageBreak bool
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

func NewProcessor(pdf Engine, doc *Document, options ...ProcessOption) (*Processor, error) {
	p := &Processor{
		doc:              doc,
		fontDir:          "fonts",
		codePage:         "",
		currStyles:       DefaultStyle,
		currPage:         0,
		preventPageBreak: false,
		pdf:              pdf,
	}
	for _, o := range options {
		err := o(p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// func (p *Processor) initFonts() error {
// 	boldOption := gofpdf.TtfOption{
// 		Style: gofpdf.Bold,
// 	}
// 	italicOption := gofpdf.TtfOption{
// 		Style: gofpdf.Italic,
// 	}
// 	boldItalicOption := gofpdf.TtfOption{
// 		Style: gofpdf.Bold | gofpdf.Italic,
// 	}

// 	//
// 	err := p.pdf.AddTTFFont("dejavusans", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("dejavusans", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", boldOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("dejavusans", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Oblique.ttf", italicOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("dejavusans", "/usr/share/fonts/truetype/dejavu/DejaVuSans-BoldOblique.ttf", boldItalicOption)
// 	if err != nil {
// 		return err
// 	}
// 	//

// 	err = p.pdf.AddTTFFont("sans", "fonts/NotoSans-Regular.ttf")
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("sans", "fonts/NotoSans-Bold.ttf", boldOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("sans", "fonts/NotoSans-Italic.ttf", italicOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("sans", "fonts/NotoSans-BoldItalic.ttf", boldItalicOption)
// 	if err != nil {
// 		return err
// 	}

// 	err = p.pdf.AddTTFFont("mono", "fonts/LiberationMono-Regular.ttf")
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("mono", "fonts/LiberationMono-Bold.ttf", boldOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("mono", "fonts/LiberationMono-Italic.ttf", italicOption)
// 	if err != nil {
// 		return err
// 	}
// 	err = p.pdf.AddTTFFontWithOption("mono", "fonts/LiberationMono-BoldItalic.ttf", boldItalicOption)
// 	if err != nil {
// 		return err
// 	}

// 	// err = p.pdf.AddTTFFont("courier", "fonts/courier.ttf")
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	err = p.pdf.AddTTFFont("wts11", "fonts/wts11.ttf")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//TODO: Units
//TODO: Orientation
//TODO: Header/Footer
//TODO: Page Numbering
//TODO: Auto Page break

// func (p *Processor) Process(w io.Writer) error {
// 	start := time.Now()
// 	// p.pdf.AliasNbPages("{np}")
// 	p.transformText = func(s string) string {
// 		ts := strings.Replace(s, "{cp}", fmt.Sprintf("%d", p.currPage), -1)
// 		return ts
// 	}

// 	fmt.Printf("run instructions ...\n")
// 	p.pageSize = gofpdf.PageSizeA4
// 	p.pdf = &gofpdf.GoPdf{}
// 	p.pdf.Start(gofpdf.Config{
// 		PageSize: *p.pageSize,
// 	})
// 	err := p.initFonts()
// 	if err != nil {
// 		return err
// 	}

// 	p.addPage()

// 	p.applyDefaults()
// 	p.applyFont(p.currStyles.Font)
// 	p.processInstructions(p.doc.Body)

// 	fmt.Printf("run instructions ... in (%s)\n", time.Since(start))
// 	return p.pdf.Write(w)
// }

func (p *Processor) Process(w io.Writer, numPages int) error {

	// p.pdf.AliasNbPages("{np}")
	p.transformText = func(s string) string {
		ts := strings.Replace(s, "{cp}", fmt.Sprintf("%d", p.currPage), -1)
		ts = strings.Replace(ts, "{np}", fmt.Sprintf("%d", numPages), -1)
		return ts
	}
	p.pageSize = gofpdf.PageSizeA4
	err := p.render()
	if err != nil {
		return err
	}

	// numPages := p.currPage
	// p.transformText = func(s string) string {
	// 	ts := strings.Replace(s, "{cp}", fmt.Sprintf("%d", p.currPage), -1)
	// 	ts = strings.Replace(ts, "{np}", fmt.Sprintf("%d", numPages), -1)
	// 	return ts
	// }
	// p.currStyles = DefaultStyle
	// p.currPage = 0
	// p.preventPageBreak = false
	// err = p.render()
	// if err != nil {
	// 	return err
	// }

	return p.pdf.Write(w)
}

func (p *Processor) render() error {
	// p.pdf = &gofpdf.GoPdf{}
	// p.pdf.Start(gofpdf.Config{
	// 	PageSize: *p.pageSize,
	// })
	c := Config{
		Format:  FormatA4,
		Margins: p.doc.Default.PageMargins,
	}
	p.pdf.Setup(c)
	// applyDefaults must apparently be called before the first add-page call
	p.applyDefaults()

	//err := p.initFonts()
	// if err != nil {
	// 	return err
	// }
	p.addPage()
	p.applyFont(p.currStyles.Font)
	p.processInstructions(p.doc.Body)
	return nil
}

func (p *Processor) addPage() {
	if p.preventPageBreak {
		return
	}
	p.pdf.AddPage()
	p.currPage++
	p.preventPageBreak = true
	p.processInstructions(p.doc.Header)
	y := p.pdf.GetY()
	p.processInstructions(p.doc.Footer)
	p.pdf.SetY(y + p.currFont.PointSize)
	p.preventPageBreak = false
}

func (p *Processor) applyDefaults() {
	//p.pdf.SetAutoPageBreak(p.doc.Default.PageBreaks == PageBreakModeAuto, p.doc.Default.PageMargins.Bottom)
	// p.pdf.SetMargins(p.doc.Default.PageMargins.Left, p.doc.Default.PageMargins.Top, p.doc.Default.PageMargins.Right, p.doc.Default.PageMargins.Bottom)
	// Logf("set margins (left, top, right, bottom) to: %f, %f, %f, %f",
	// 	p.doc.Default.PageMargins.Left, p.doc.Default.PageMargins.Top, p.doc.Default.PageMargins.Right, p.doc.Default.PageMargins.Bottom)
	//TODO: included in setup
}

func (p *Processor) appliedStyles(i Instruction) style.Styles {
	st := p.currStyles
	i.Apply(p.doc.styleClasses, &st)
	return st
}

func dumpStyles(st style.Styles) {
	b, _ := json.MarshalIndent(st, "", "  ")
	Logf("styles:\n%s\n", string(b))
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
			if !i.FromBottom {
				p.pdf.SetY(i.Y)
			} else {
				p.pdf.SetY(p.pageSize.H - i.Y)
			}
		case *SetXY:
			p.pdf.SetX(i.X)
			p.pdf.SetY(i.Y)
		case *Box:
			p.renderTextBox(i.Text, p.appliedStyles(i))
		case *Text:
			p.renderText(i, p.appliedStyles(i))
		case *Table:
			//fmt.Printf("\n***\ntable: %#v\n***\n", i)
			c := i.Clone()
			st := p.appliedStyles(c)
			//c.Dump()
			//dumpStyles(st)
			p.renderTable(c, st)
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
	//l, _, r, _ := p.pdf.Margins()
	ms := p.pdf.Margins()
	return pw - (ms.Left + ms.Right) - 3 // without substracting 3 it doesn't fit
}

func (p *Processor) applyFont(fnt style.Font) {
	if fnt == p.currFont {
		return
	}
	// st := fpdfFontStyle(fnt)
	// //Logf("apply-font: %s -> %q", fnt, st)
	// p.currFont = fnt
	// p.pdf.SetFont(string(fnt.Family), st, int(fnt.PointSize))
	p.currFont = fnt
	p.pdf.UseFont(fnt)
}

func (p *Processor) processLineFeed(lf *LineFeed) {
	//_, fontHeight := p.pdf.GetFontSize()
	fontHeight := p.currFont.PointSize
	height := fontHeight * lf.Lines
	p.pdf.SetY(p.pdf.GetY() + height)
}

func (p *Processor) ln(h float64) {
	ms := p.pdf.Margins()
	p.pdf.SetX(ms.Left)
	y := p.pdf.GetY()
	if y+h >= p.pageSize.H {
		p.addPage()
	} else {
		p.pdf.SetY(y + h)
	}
}

func (p *Processor) renderText(text *Text, sty style.Styles) {
	p.write(text.Text, p.effectiveWidth(sty.Dimension.Width), sty.Dimension.LineHeight, sty.Align.HAlign, sty.Font, sty.Color.Text)
}

func (p *Processor) GetXY() (float64, float64) {
	//return p.pdf.GetX() + p.pdf.MarginLeft(), p.pdf.GetY()
	return p.pdf.GetX(), p.pdf.GetY()
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
	ph := p.pageSize.H

	if y0+height >= ph {
		p.addPage()
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
	p.ln(sty.Dimension.LineHeight + sty.Box.Padding.Bottom)
}

func (p *Processor) renderImage(img *Image, sty style.Styles) {
	x0, y0 := p.GetXY()
	x0 += sty.Dimension.OffsetX
	y0 += sty.Dimension.OffsetY
	// r := gofpdf.Rect{
	// 	W: sty.Dimension.Width,
	// 	H: sty.Dimension.Height,
	// }
	//err := p.pdf.Image(img.Source, x0, y0, &r)
	err := p.pdf.Image(img.Source, x0, y0, sty.Dimension.Width, sty.Dimension.Height)
	if err != nil {
		Logf("image: %v", err)
	}
}
