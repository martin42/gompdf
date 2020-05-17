package gompdf

import (
	"encoding/json"
	"encoding/xml"

	"github.com/mazzegi/gompdf/style"
	"github.com/pkg/errors"
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

func (cell *TableCell) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				Logf("decode cell instruction failed: %v", err)
				continue
			}
			cell.Instructions = append(cell.Instructions, i)
		case xml.CharData:
			cell.Content += string(t)
		}
	}
}

type Table struct {
	Styled
	XMLName xml.Name    `xml:"table"`
	Rows    []*TableRow `xml:"tr"`
}

type TableRow struct {
	Styled
	XMLName xml.Name     `xml:"tr"`
	Cells   []*TableCell `xml:"td"`
}

type TableCell struct {
	Styled
	XMLName        xml.Name `xml:"td"`
	Content        string   `xml:",chardata"`
	Instructions   []Instruction
	spannedBy      *TableCell
	spans          []*TableCell
	x0, y0, x1, y1 float64
}

func (t *Table) Dump() {
	b, _ := json.MarshalIndent(t, "", "  ")
	Logf("table:\n%s\n", string(b))
}

func (t *Table) Clone() *Table {
	ct := &Table{
		Styled: t.Styled,
	}
	for _, r := range t.Rows {
		ct.Rows = append(ct.Rows, r.Clone())
	}
	return ct
}

func (r *TableRow) Clone() *TableRow {
	cr := &TableRow{
		Styled: r.Styled,
	}
	//Logf("clone-row with %d cells", len(r.Cells))
	for _, c := range r.Cells {
		cr.Cells = append(cr.Cells, c.Clone())
	}
	return cr
}

func (c *TableCell) Clone() *TableCell {
	cc := &TableCell{
		Styled:       c.Styled,
		Content:      c.Content,
		Instructions: c.Instructions,
	}
	// for _, i := range c.Instructions {
	// 	cc.Instructions = append(cc.Instructions, i)
	// }
	return cc
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
			if cellStyles.Table.ColumnWidth > 0 {
				cw = cellStyles.Table.ColumnWidth
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

func (p *Processor) processTableSpans(t *Table) error {
	for ir, row := range t.Rows {
		for ic, cell := range row.Cells {
			var cellStyles style.Styles
			cell.Apply(style.Classes{}, &cellStyles)
			if cellStyles.Table.RowSpan > 1 {
				//Logf("found row-span (%d) (%d) -> (%d)", ir, ic, cellStyles.Table.RowSpan)
				//insert spanned cell in following rows
				for n := 0; n <= cellStyles.RowSpan-2; n++ {
					spannedRowIdx := ir + 1 + n
					if spannedRowIdx >= len(t.Rows) {
						return errors.Errorf("row span exceeds table")
					}
					spannedRow := t.Rows[spannedRowIdx]
					//Logf("spanned row has %d cells", len(spannedRow.Cells))
					if ic <= len(spannedRow.Cells) {
						newCells := []*TableCell{}
						for _, c := range spannedRow.Cells[:ic] {
							newCells = append(newCells, c)
						}
						newCell := &TableCell{
							spannedBy: cell,
						}
						cell.spans = append(cell.spans, newCell)
						newCells = append(newCells, newCell)
						newCells = append(newCells, spannedRow.Cells[ic:]...)
						spannedRow.Cells = newCells
						//Logf("new row (%d): has %d cells", spannedRowIdx, len(spannedRow.Cells))
					}
				}
			}
		}
	}
	return nil
}

func (p *Processor) tableHeight(t *Table, tableStyles style.Styles) float64 {
	if t.MaxColumnCount() == 0 {
		return 0
	}
	cellHeight := func(c *TableCell, cellWidth float64, cellStyle style.Styles) float64 {
		textWidth := cellWidth - cellStyle.Box.Padding.Left - cellStyle.Box.Padding.Right
		height := p.textHeight(c.Content, textWidth, cellStyle.Dimension.LineHeight, cellStyle.Font)
		return height + cellStyle.Box.Padding.Top + cellStyle.Box.Padding.Bottom
	}

	//widthTotal, _ := p.pdf.GetPageSize()
	widthTotal := p.pageSize.W
	leftM, _, rightM, _ := p.pdf.Margins()
	widthTotal -= (leftM + rightM)
	colWs := p.ColumnWidths(t, widthTotal, tableStyles)

	totalHeight := float64(0)
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
		totalHeight += rowHeight
	}
	return totalHeight
}

func (p *Processor) renderTable(t *Table, tableStyles style.Styles) {
	if t.MaxColumnCount() == 0 {
		return
	}
	// in order to remove temporary span variables/values
	//defer t.resetSpans()

	p.processTableSpans(t)
	tableHeight := p.tableHeight(t, tableStyles)

	//if not further specified, distribute witdths uniformly
	widthTotal := p.pageSize.W
	leftM, _, rightM, bottomM := p.pdf.Margins()
	widthTotal -= (leftM + rightM)
	colWs := p.ColumnWidths(t, widthTotal, tableStyles)
	//Logf("col-widths: %v", colWs)

	cellHeight := func(c *TableCell, cellIdx int, cellStyle style.Styles) float64 {
		cellWidth := colWs[cellIdx]
		if cellStyle.Table.ColumnSpan > 1 {
			for cs := 1; cs < cellStyle.Table.ColumnSpan; cs++ {
				if cellIdx+cs < len(colWs) {
					cellWidth += colWs[cellIdx+cs]
				}
			}
		}
		textWidth := cellWidth - cellStyle.Box.Padding.Left - cellStyle.Box.Padding.Right
		height := p.textHeight(c.Content, textWidth, cellStyle.Dimension.LineHeight, cellStyle.Font)
		return height + cellStyle.Box.Padding.Top + cellStyle.Box.Padding.Bottom
	}

	ph := p.pageSize.H
	ph -= (bottomM)
	x0 := p.pdf.GetX()
	y := p.pdf.GetY()
	if y+tableHeight > ph {
		p.addPage()
		y = p.pdf.GetY()
	}

	afterRender := []func(){}

	for ir, row := range t.Rows {
		rowStyles := tableStyles
		if ir == 0 {
			row.ApplyWithSelector("first", p.doc.styleClasses, &rowStyles)
		} else if ir == len(t.Rows)-1 {
			row.ApplyWithSelector("last", p.doc.styleClasses, &rowStyles)
		} else {
			row.Apply(p.doc.styleClasses, &rowStyles)
		}
		//calc row height
		rowHeight := float64(0)
		for i, c := range row.Cells {
			if len(c.spans) > 0 || c.spannedBy != nil {
				//don't consider row-spanned cells for general rowheight, as their height might by bigger
				continue
			}
			cellStyles := rowStyles
			c.Apply(p.doc.styleClasses, &cellStyles)
			ch := cellHeight(c, i, cellStyles)
			if ch > rowHeight {
				rowHeight = ch
			}
		}

		if y+rowHeight >= ph {
			p.addPage()
			y = p.pdf.GetY()
		}

		x := x0
		colOffset := 0
		for ic, c := range row.Cells {
			cellStyles := rowStyles
			if ic == 0 {
				c.ApplyWithSelector("first", p.doc.styleClasses, &cellStyles)
			} else if ic == len(row.Cells) {
				c.ApplyWithSelector("last", p.doc.styleClasses, &cellStyles)
			} else {
				c.Apply(p.doc.styleClasses, &cellStyles)
			}

			p.SetXY(x, y)
			x0 := x
			y0 := y
			ws := colWs[colOffset]
			x1 := x + ws
			if cellStyles.Table.ColumnSpan > 1 {
				for cs := 1; cs < cellStyles.Table.ColumnSpan; cs++ {
					if colOffset+cs < len(colWs) {
						x1 += colWs[colOffset+cs]
						ws += colWs[colOffset+cs]
					}
				}
			}
			colOffset += cellStyles.Table.ColumnSpan

			y1 := y + rowHeight
			c.x0, c.y0, c.x1, c.y1 = x0, y0, x1, y1
			x += ws
			if c.spannedBy != nil {
				continue
			}
			if len(c.spans) > 0 {
				arStyles := cellStyles
				arCell := c
				afterRender = append(afterRender, func() {
					//collect bounds
					y1 := arCell.y1
					for _, sc := range arCell.spans {
						if sc.y1 > y1 {
							y1 = sc.y1
						}
					}
					p.renderCell(arCell.x0, arCell.y0, arCell.x1, y1, arCell, arStyles)
				})
				continue
			}
			p.renderCell(x0, y0, x1, y1, c, cellStyles)
		}
		y += rowHeight
		p.ln(-1)
	}

	for _, ar := range afterRender {
		ar()
	}

	p.SetXY(x0, y)
	p.ln(-1)
}

func (p *Processor) renderCell(x0, y0, x1, y1 float64, c *TableCell, cellStyles style.Styles) {
	p.drawBox(x0, y0, x1, y1, cellStyles)

	textWidth := (x1 - x0) - cellStyles.Box.Padding.Left - cellStyles.Box.Padding.Right - 3 //wihout 2 it doesn't fit
	textHeight := p.textHeight(c.Content, textWidth, cellStyles.Dimension.LineHeight, cellStyles.Font)
	textMargin := y1 - y0 - textHeight - cellStyles.Box.Padding.Top - cellStyles.Box.Padding.Bottom
	if textMargin < 0 {
		textMargin = 0
	}

	if cellStyles.Align.VAlign == style.VAlignMiddle {
		p.pdf.SetY(y0 + cellStyles.Box.Padding.Top + textMargin/2)
	} else if cellStyles.Align.VAlign == style.VAlignBottom {
		p.pdf.SetY(y0 + cellStyles.Box.Padding.Top + textMargin)
	} else {
		p.pdf.SetY(y0 + cellStyles.Box.Padding.Top)
	}
	p.pdf.SetX(x0 + cellStyles.Box.Padding.Left)

	if c.Content != "" {
		p.write(c.Content, textWidth, cellStyles.Dimension.LineHeight, cellStyles.Align.HAlign, cellStyles.Font, cellStyles.Color.Text)
	}

	for _, inst := range c.Instructions {
		switch inst := inst.(type) {
		case *Box:
			p.pdf.SetY(y0 + cellStyles.Box.Padding.Top)
			p.pdf.SetX(x0 + cellStyles.Box.Padding.Left)
			p.renderTextBox(inst.Text, p.appliedStyles(inst))
		case *Image:
			p.pdf.SetY(y0 + cellStyles.Box.Padding.Top)
			p.pdf.SetX(x0 + cellStyles.Box.Padding.Left)
			p.renderImage(inst, p.appliedStyles(inst))
		}
	}
}
