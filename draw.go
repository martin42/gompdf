package gompdf

import "github.com/mazzegi/gompdf/style"

func (p *Processor) drawRect(x0, y0, x1, y1 float64) {
	p.pdf.Line(x0, y0, x1, y0)
	p.pdf.Line(x1, y0, x1, y1)
	p.pdf.Line(x1, y1, x0, y1)
	p.pdf.Line(x0, y1, x0, y0)
}

func (p *Processor) drawBox(x0, y0, x1, y1 float64, sty style.Styles) {
	p.pdf.SetLineWidth(sty.Draw.LineWidth)
	p.pdf.SetStrokeColor(sty.Color.Foreground)
	p.pdf.SetFillColor(sty.Color.Foreground)

	width := x1 - x0
	height := y1 - y0
	p.pdf.FillRect(x0, y0, width, height)
	//p.pdf.MoveTo(x0-sty.Draw.LineWidth/2, y0)
	if sty.Box.Border.Top > 0 {
		p.pdf.Line(x0, y0, x1, y0)
		//p.pdf.LineTo(x0+width, y0)
	} else {
		//p.pdf.MoveTo(x0+width, y0)
	}
	if sty.Box.Border.Right > 0 {
		//p.pdf.LineTo(x0+width, y1)
		p.pdf.Line(x1, y0, x1, y1)
	} else {
		//p.pdf.MoveTo(x0+width, y1)
	}
	if sty.Box.Border.Bottom > 0 {
		//p.pdf.LineTo(x0, y1)
		p.pdf.Line(x1, y1, x0, y1)
	} else {
		//p.pdf.MoveTo(x0, y1)
	}
	if sty.Box.Border.Left > 0 {
		//p.pdf.LineTo(x0, y0)
		p.pdf.Line(x0, y1, x0, y0)
	} else {
		//p.pdf.MoveTo(x0, y0)
	}
	//p.pdf.DrawPath("D")
}
