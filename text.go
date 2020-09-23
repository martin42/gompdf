package gompdf

import (
	"strings"

	"github.com/mazzegi/gompdf/markdown"
	"github.com/mazzegi/gompdf/style"
)

func (p *Processor) normalizedText(s string) string {
	//remove carriage return and tabs
	text := s
	text = strings.Replace(text, "\r", "\n", -1)
	text = strings.Replace(text, "\t", " ", -1)
	//split into lines
	lines := strings.Split(text, "\n")
	tlines := []string{}
	for _, line := range lines {
		if line != "" {
			tlines = append(tlines, strings.Trim(line, " "))
		}
	}
	sn := strings.Join(tlines, " ")
	var snt string
	for {
		snt = strings.Replace(sn, "  ", " ", -1)
		if snt == sn {
			return snt
		}
		sn = snt
	}
}

type textLine struct {
	mdWords               markdown.Items
	textWidth             float64
	textWidthTrimmedRight float64
}

func (p *Processor) applyMarkdownFont(mdi markdown.Item, toFnt style.Font) {
	mdFont := toFnt
	fntStyles := fpdfFontStyle(toFnt)
	if mdi.Italic {
		fntStyles += "I"
		mdFont.Style = style.FontStyleItalic
	}
	if mdi.Bold {
		fntStyles += "B"
		mdFont.Weight = style.FontWeightBold
	}
	family := string(toFnt.Family)
	if mdi.Code {
		family = "mono"
		mdFont.Family = family
	}
	//p.pdf.SetFont(family, fntStyles, toFnt.PointSize)
	p.applyFont(mdFont)
}

func (p *Processor) textLines(mdWords markdown.Items, width float64, fnt style.Font) []textLine {
	lines := []textLine{}
	currLine := textLine{
		mdWords:   markdown.Items{},
		textWidth: 0.0,
	}
	for _, mdWord := range mdWords {
		if mdWord.Newline {
			lines = append(lines, currLine)
			currLine = textLine{
				mdWords:   markdown.Items{},
				textWidth: 0,
			}
			continue
		}
		p.applyMarkdownFont(mdWord, fnt)
		// wordWidth, _ := p.pdf.MeasureTextWidth(mdWord.Text)
		// wordWidthTrimmedRight, _ := p.pdf.MeasureTextWidth(strings.TrimRight(mdWord.Text, " "))
		wordWidth := p.pdf.TextWidth(mdWord.Text)
		wordWidthTrimmedRight := p.pdf.TextWidth(strings.TrimRight(mdWord.Text, " "))
		Logf("measure %q: w=%f (trimmed=%f)", mdWord.Text, wordWidth, wordWidthTrimmedRight)
		if currLine.textWidth+wordWidth > width {
			lines = append(lines, currLine)
			currLine = textLine{
				mdWords:   markdown.Items{},
				textWidth: 0,
			}
		}
		if len(currLine.mdWords) == 0 {
			mdWord.Text = strings.TrimLeft(mdWord.Text, " ")
			//wordWidth, _ = p.pdf.MeasureTextWidth(mdWord.Text)
			wordWidth = p.pdf.TextWidth(mdWord.Text)
		}
		currLine.mdWords = append(currLine.mdWords, mdWord)
		currLine.textWidthTrimmedRight = currLine.textWidth + wordWidthTrimmedRight
		currLine.textWidth += wordWidth
	}
	if len(currLine.mdWords) > 0 {
		lines = append(lines, currLine)
	}
	for _, l := range lines {
		if len(l.mdWords) > 0 {
			l.mdWords[0].Text = strings.TrimLeft(l.mdWords[0].Text, " ")
			l.mdWords[len(l.mdWords)-1].Text = strings.TrimRight(l.mdWords[len(l.mdWords)-1].Text, " ")
		}
	}
	return lines
}

func (p *Processor) write(text string, width float64, lineHeight float64, halign style.HAlign, fnt style.Font, cr style.RGB) {
	p.applyFont(fnt)
	//p.pdf.SetTextColor(cr.R, cr.G, cr.B)
	p.pdf.SetTextColor(cr)
	text = p.normalizedText(text)
	fontHeight := p.currFont.PointSize
	height := fontHeight * lineHeight
	xLeft := p.pdf.GetX()
	p.ln(fontHeight)
	mdWords := markdown.NewProcessor().Process(text).WordItems(p.transformText)
	lines := p.textLines(mdWords, width, fnt)
	for _, line := range lines {
		if len(line.mdWords) == 0 {
			continue
		}
		switch halign {
		case style.HAlignLeft:
			p.pdf.SetX(xLeft)
		case style.HAlignCenter:
			p.pdf.SetX(xLeft + (width-line.textWidthTrimmedRight)/2.0)
		case style.HAlignRight:
			p.pdf.SetX(xLeft + width - line.textWidthTrimmedRight)
		}

		for _, mdWord := range line.mdWords {
			if mdWord.Text == "" {
				continue
			}
			p.applyMarkdownFont(mdWord, fnt)
			//p.pdf.Write(height, mdWord.Text)
			Logf("write %q (x == %f)", mdWord.Text, p.pdf.GetX())
			p.pdf.Text(mdWord.Text)
		}
		p.ln(height)
	}
	//p.pdf.SetTextColor(p.currStyles.Color.Text.R, p.currStyles.Color.Text.G, p.currStyles.Color.Text.B)
	p.pdf.SetTextColor(p.currStyles.Color.Text)
	p.applyFont(p.currStyles.Font)
}

func (p *Processor) textHeight(text string, width float64, lineHeight float64, fnt style.Font) float64 {
	p.applyFont(fnt)
	text = p.normalizedText(text)
	fontHeight := p.currFont.PointSize
	height := fontHeight * lineHeight
	mdWords := markdown.NewProcessor().Process(text).WordItems(p.transformText)
	lines := p.textLines(mdWords, width, fnt)
	textHeight := float64(0)
	for _, line := range lines {
		if len(line.mdWords) == 0 {
			continue
		}
		textHeight += height
	}
	p.applyFont(p.currStyles.Font)
	return textHeight
}
