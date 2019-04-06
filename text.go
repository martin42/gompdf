package gompdf

import (
	"strings"

	"github.com/martin42/gompdf/markdown"
)

type textLine struct {
	mdWords               markdown.Items
	textWidth             float64
	textWidthTrimmedRight float64
}

func (p *Processor) applyMarkdownFont(mdi markdown.Item) {
	fntStyles := ""
	if mdi.Italic {
		fntStyles += "I"
	}
	if mdi.Bold {
		fntStyles += "B"
	}
	family := string(p.currFont.fontFamily)
	if mdi.Code {
		family = "Courier"
	}
	p.pdf.SetFont(family, fntStyles, float64(p.currFont.fontPointSize))
}

func (p *Processor) textLines(mdWords markdown.Items, width float64) []textLine {
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
		p.applyMarkdownFont(mdWord)
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

func (p *Processor) write(text string, width float64, lineHeight float64, halign HAlign) {
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Trim(text, " ")
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lineHeight
	xLeft := p.pdf.GetX()
	mdWords := markdown.NewProcessor().Process(text).WordItems()
	lines := p.textLines(mdWords, width)
	for _, line := range lines {
		if len(line.mdWords) == 0 {
			continue
		}
		switch halign {
		case HAlignLeft:
			p.pdf.SetX(xLeft)
		case HAlignCenter:
			p.pdf.SetX(xLeft + (width-line.textWidthTrimmedRight)/2.0)
		case HAlignRight:
			p.pdf.SetX(xLeft + width - line.textWidthTrimmedRight)
		}

		for _, mdWord := range line.mdWords {
			p.applyMarkdownFont(mdWord)
			p.pdf.Write(height, mdWord.Text)
		}
		p.pdf.Ln(height)
	}
	p.applyFont(&p.currFont)
}

func (p *Processor) textHeight(text string, width float64, lineHeight float64) float64 {
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Trim(text, " ")
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lineHeight
	mdWords := markdown.NewProcessor().Process(text).WordItems()
	lines := p.textLines(mdWords, width)
	textHeight := float64(0)
	for _, line := range lines {
		if len(line.mdWords) == 0 {
			continue
		}
		textHeight += height
	}
	p.applyFont(&p.currFont)
	return textHeight
}
