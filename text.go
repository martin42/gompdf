package gompdf

import (
	"strings"

	"github.com/martin42/gompdf/markdown"
)

type textLine struct {
	mdWords   markdown.Items
	textWidth float64
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
		if currLine.textWidth+wordWidth > width {
			lines = append(lines, currLine)
			mdWord.Text = strings.TrimLeft(mdWord.Text, " ")
			currLine = textLine{
				mdWords:   markdown.Items{mdWord},
				textWidth: wordWidth,
			}
		} else {
			if len(currLine.mdWords) == 0 {
				mdWord.Text = strings.TrimLeft(mdWord.Text, " ")
			}
			currLine.mdWords = append(currLine.mdWords, mdWord)
			currLine.textWidth += wordWidth
		}
	}
	if len(currLine.mdWords) > 0 {
		lines = append(lines, currLine)
	}
	return lines
}

func (p *Processor) write(text string, width float64, lineHeight float64, halign HAlign) {
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "{nl}", "\n", -1)
	text = strings.Trim(text, " ")
	_, fontHeight := p.pdf.GetFontSize()
	height := fontHeight * lineHeight
	xLeft := p.pdf.GetX()
	mdWords := markdown.NewProcessor().Process(text).WordItems()
	lines := p.textLines(mdWords, width)
	for _, line := range lines {
		switch halign {
		case HAlignLeft:
			p.pdf.SetX(xLeft)
		case HAlignCenter:
			p.pdf.SetX(xLeft + (width-line.textWidth)/2.0)
		case HAlignRight:
			p.pdf.SetX(xLeft + width - line.textWidth)
		}

		for _, mdWord := range line.mdWords {
			p.applyMarkdownFont(mdWord)
			p.pdf.Write(height, mdWord.Text)
		}
		p.pdf.Ln(height)
	}
	p.processFont(&p.currFont)
}
