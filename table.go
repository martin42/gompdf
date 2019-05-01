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

func (p *Processor) ColumnWidths(t *Table, pageWidth float64, tableStyles style.Styles) []float64 {
	cws := make([]float64, t.MaxColumnCount())
	for _, row := range t.Rows {
		rowStyles := tableStyles
		row.Apply(p.doc.styleClasses, &rowStyles)
		for i, c := range row.Cells {
			cellStyles := rowStyles
			c.Apply(p.doc.styleClasses, &cellStyles)
			cw := float64(-1)
			if cellStyles.Dimension.ColumnWidth > 0 {
				cw = cellStyles.Dimension.ColumnWidth
			}
			if cw > 0 && cw > cws[i] {
				cws[i] = cw
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

func (p *Processor) renderTable(t *Table, tableStyles style.Styles) {
	if t.MaxColumnCount() == 0 {
		return
	}

	cellHeight := func(c *TableCell, cellWidth float64, cellStyle style.Styles) float64 {
		textWidth := cellWidth - cellStyle.Box.Padding.Left - cellStyle.Box.Padding.Right
		height := p.textHeight(c.Content, textWidth, cellStyle.Dimension.LineHeight, cellStyle.Font)
		return height + cellStyle.Box.Padding.Top + cellStyle.Box.Padding.Bottom
	}

	//if not further specified, distribute witdths uniformly
	widthTotal, _ := p.pdf.GetPageSize()
	leftM, _, rightM, _ := p.pdf.GetMargins()
	widthTotal -= (leftM + rightM)
	colWs := p.ColumnWidths(t, widthTotal, tableStyles)

	x0 := p.pdf.GetX()
	y := p.pdf.GetY()
	for _, row := range t.Rows {
		rowStyles := tableStyles
		row.Apply(p.doc.styleClasses, &rowStyles)
		//calc row height
		rowHeight := float64(0)
		for i, c := range row.Cells {
			cellStyles := rowStyles
			c.Apply(p.doc.styleClasses, &cellStyles)
			ch := cellHeight(c, colWs[i], cellStyles)
			if ch > rowHeight {
				rowHeight = ch
			}
		}
		Logf("row-height: %.1f", rowHeight)
		x := x0
		for i, c := range row.Cells {
			cellStyles := rowStyles
			c.Apply(p.doc.styleClasses, &cellStyles)
			p.pdf.SetXY(x, y)

			x0 := x
			y0 := y
			x1 := x + colWs[i]
			y1 := y + rowHeight
			p.drawBox(x0, y0, x1, y1, cellStyles)

			//Reset, to start writing at top left
			p.pdf.SetY(y0 + cellStyles.Box.Padding.Top)
			p.pdf.SetX(x0 + cellStyles.Box.Padding.Left)

			textWidth := colWs[i] - cellStyles.Box.Padding.Left - cellStyles.Box.Padding.Right
			p.write(c.Content, textWidth, cellStyles.Dimension.LineHeight, cellStyles.Align.HAlign, cellStyles.Font)

			x += colWs[i]
		}
		y += rowHeight
		p.pdf.Ln(-1)
	}
}
