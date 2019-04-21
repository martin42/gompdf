package gompdf

import (
	"encoding/xml"

	"github.com/martin42/gompdf/style"
)

func (tab *Table) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t == start.End() {
				return nil
			}
		case xml.StartElement:
			i, err := instructionRegistry.Decode(d, t)
			if err != nil {
				return err
			}
			switch i := i.(type) {
			case *TableRow:
				tab.Rows = append(tab.Rows, i)
			}
		}
	}
}

func (row *TableRow) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t == start.End() {
				return nil
			}
		case xml.StartElement:
			i, err := instructionRegistry.Decode(d, t)
			if err != nil {
				return err
			}
			switch i := i.(type) {
			case *TableCell:
				row.Cells = append(row.Cells, i)
			}
		}
	}
}

type Table struct {
	Styled
	XMLName xml.Name    `xml:"Table"`
	Rows    []*TableRow `xml:"Tr"`
}

type TableRow struct {
	Styled
	XMLName xml.Name     `xml:"Tr"`
	Cells   []*TableCell `xml:"Td"`
}

type TableCell struct {
	Styled
	XMLName xml.Name `xml:"Td"`
	Content string   `xml:",chardata"`
}

func (p *Processor) ColumnWidths(t *Table, pageWidth float64, sty style.Styles) []float64 {
	cellWidth := func(c *TableCell, def float64) float64 {
		s := p.appliedStyles(c)
		if s.Dimension.ColumnWidth > 0 {
			return s.Dimension.ColumnWidth
		}
		return def
	}

	cws := make([]float64, t.MaxColumnCount())
	for _, row := range t.Rows {
		for i, c := range row.Cells {
			w := cellWidth(c, -1)
			if w > 0 && w > cws[i] {
				cws[i] = w
			}
		}
	}
	spaceUsed := float64(0)
	columnsZero := float64(0)
	for _, c := range cws {
		if c > 0 {
			spaceUsed += c
		} else {
			columnsZero++
		}
	}
	if columnsZero > 0 {
		colWidth := (pageWidth - spaceUsed) / columnsZero
		for i, c := range cws {
			if c == 0 {
				cws[i] = colWidth
			}
		}
	}
	return cws
}

func (t Table) MaxColumnCount() int {
	m := 0
	for _, row := range t.Rows {
		if len(row.Cells) > m {
			m = len(row.Cells)
		}
	}
	return m
}

func (p *Processor) renderTable(t *Table, sty style.Styles) {
	if t.MaxColumnCount() == 0 {
		return
	}

	cellHeight := func(c *TableCell, cellWidth float64) float64 {
		s := p.appliedStyles(c)
		textWidth := cellWidth - s.Box.Padding.Left - s.Box.Padding.Right
		height := p.textHeight(c.Content, textWidth, s.Dimension.LineHeight, s.Font)
		return height + s.Box.Padding.Top + s.Box.Padding.Bottom
	}

	//if not further specified, distribute witdths uniformly
	widthTotal, _ := p.pdf.GetPageSize()
	leftM, _, rightM, _ := p.pdf.GetMargins()
	widthTotal -= (leftM + rightM)
	colWs := p.ColumnWidths(t, widthTotal, sty)

	x0 := p.pdf.GetX()
	y := p.pdf.GetY()
	for _, row := range t.Rows {
		//calc row height
		rowHeight := float64(0)
		for i, c := range row.Cells {
			ch := cellHeight(c, colWs[i])
			if ch > rowHeight {
				rowHeight = ch
			}
		}

		x := x0
		for i, c := range row.Cells {
			Logf("cell (%d): (x=%.1f) (y=%.1f)", i, x, y)
			p.pdf.SetXY(x, y)

			s := p.appliedStyles(c)
			x0 := x
			y0 := y
			x1 := x + colWs[i]
			y1 := y + rowHeight
			p.drawBox(x0, y0, x1, y1, s)

			//Reset, to start writing at top left
			p.pdf.SetY(y0 + s.Box.Padding.Top)
			p.pdf.SetX(x0 + s.Box.Padding.Left)

			textWidth := colWs[i] - s.Box.Padding.Left - s.Box.Padding.Right
			p.write(c.Content, textWidth, s.Dimension.LineHeight, s.Align.HAlign, s.Font)

			x += colWs[i]
		}
		y += rowHeight
		p.pdf.Ln(-1)
	}
}
