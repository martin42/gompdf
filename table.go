package gompdf

import (
	"encoding/xml"
)

type Table struct {
	NoStyles
	XMLName xml.Name   `xml:"Table"`
	Rows    []TableRow `xml:"Tr"`
}

type TableRow struct {
	XMLName xml.Name       `xml:"Tr"`
	Cells   []Instructions `xml:"Td"`
}

func (t Table) ColumnWidths(pageWidth float64) []float64 {
	cellWidth := func(c Instructions, def float64) float64 {
		for _, s := range c.Styles {
			switch i := s.(type) {
			case Width:
				return float64(i)
			}
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

func (p *Processor) renderTable(t *Table) {
	if t.MaxColumnCount() == 0 {
		return
	}

	//if not further specified, distribute witdths uniformly
	widthTotal, _ := p.pdf.GetPageSize()
	leftM, _, rightM, _ := p.pdf.GetMargins()
	widthTotal -= (leftM + rightM)
	colWs := t.ColumnWidths(widthTotal)

	x0 := p.pdf.GetX()
	for _, row := range t.Rows {
		x := x0
		y := p.pdf.GetY()
		for i, c := range row.Cells {
			Logf("cell (%d): (x=%.1f) (y=%.1f)", i, x, y)
			p.pdf.SetXY(x, y)
			for _, is := range c.iss {
				switch is := is.(type) {
				case *Text:
					p.write(is.Text, colWs[i], 1.5, HAlignLeft)
				}
			}
			x += colWs[i]
		}
		p.pdf.Ln(-1)
	}
}
