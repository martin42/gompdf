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
