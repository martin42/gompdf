package gompdf

import (
	"strings"

	"github.com/martin42/gompdf/markdown"
	"github.com/martin42/gompdf/style"
)

type textLine struct {
	mdWords               markdown.Items
	textWidth             float64
	textWidthTrimmedRight float64
}

func (p *Processor) applyMarkdownFont(mdi markdown.Item, toFnt style.Font) {
	fntStyles := fpdfFontStyle(toFnt) // ""
	if mdi.Italic {
		fntStyles += "I"
	}
	if mdi.Bold {
		fntStyles += "B"
	}
	family := string(toFnt.Family)
	if mdi.Code {
		family = "Courier"
	}
	p.pdf.SetFont(family, fntStyles, toFnt.PointSize)
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
		wordWidth := p.pdf.GetStringWidth(mdWord.Text)
		wordWidthTrimmedRight := p.pdf.GetStringWidth(strings.TrimRight(mdWord.Text, " "))
		if currLine.textWidth+wordWidth > width {
			lines = append(lines, currLine)
			currLine = textLine{
				mdWords:   markdown.Items{},
				textWidth: 0,
			}
		}
		if len(currLine.mdWords) == 0 {
			mdWord.Text = strings.TrimLeft(mdWord.Text, " ")
			wordWidth = p.pdf.GetStringWidth(mdWord.Text)
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
			l.mdWords[len(l.mdWords)-1].Text = strings.TrimRight(l.mdWords[len(l.mdWords)-1].Text, " ")
		}
	}
	return lines
}

func (p *Processor) write(text string, width float64, lineHeight float64, halign style.HAlign, fnt style.Font) {
	Logf("write: font-weight: %s", fnt.Weight)

	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Trim(text, " ")
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lineHeight
	xLeft := p.pdf.GetX()
	mdWords := markdown.NewProcessor().Process(text).WordItems()
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
			p.applyMarkdownFont(mdWord, fnt)
			p.pdf.Write(height, mdWord.Text)
		}
		p.pdf.Ln(height)
	}
	p.applyFont(fnt)
}

func (p *Processor) textHeight(text string, width float64, lineHeight float64, fnt style.Font) float64 {
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Trim(text, " ")
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lineHeight
	mdWords := markdown.NewProcessor().Process(text).WordItems()
	lines := p.textLines(mdWords, width, fnt)
	textHeight := float64(0)
	for _, line := range lines {
		if len(line.mdWords) == 0 {
			continue
		}
		textHeight += height
	}
	p.applyFont(fnt)
	return textHeight
}
