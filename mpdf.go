package gompdf

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/pkg/errors"
)

func ParseAndBuild(source string, target string) error {
	doc, err := LoadFromFile(source)
	if err != nil {
		return err
	}
	outF, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outF.Close()
	p, err := NewProcessor(doc)
	if err != nil {
		return err
	}
	err = p.Process(outF)
	if err != nil {
		return err
	}
	return nil
}

func Load(r io.Reader) (*Document, error) {
	doc := &Document{}
	err := xml.NewDecoder(r).Decode(doc)
	if err != nil {
		return nil, err
	}
	err = doc.applyClasses(&doc.Body)
	if err != nil {
		return nil, err
	}
	err = doc.applyClasses(&doc.Header)
	if err != nil {
		return nil, err
	}
	err = doc.applyClasses(&doc.Footer)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func LoadFromFile(file string) (*Document, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Errorf("open (%s)", file)
	}
	defer f.Close()
	return Load(f)
}

func (doc *Document) applyClasses(iss *Instructions) error {
	for _, is := range *iss {
		is.ApplyClasses(doc.Style)
	}
	return nil
}

type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

type Unit string

const (
	UnitPt   Unit = "pt"
	UnitMm   Unit = "mm"
	UnitCm   Unit = "cm"
	UnitInch Unit = "in"
)

type Format string

const (
	FormatA3     Format = "a3"
	FormatA4     Format = "a4"
	FormatA5     Format = "a5"
	FormatLetter Format = "letter"
	FormatLegal  Format = "legal"
)

type PageBreakMode string

const (
	PageBreakModeAuto   PageBreakMode = "auto"
	PageBreakModeManual PageBreakMode = "manual"
)

type PageMargins struct {
	XMLName xml.Name `xml:"PageMargins"`
	Left    float64  `xml:"Left"`
	Top     float64  `xml:"Top"`
	Right   float64  `xml:"Right"`
	Bottom  float64  `xml:"Bottom"`
}

type Document struct {
	XMLName xml.Name     `xml:"Document"`
	Meta    Meta         `xml:"Meta"`
	Default Default      `xml:"Default"`
	Style   StyleClasses `xml:"Style"`
	Header  Instructions `xml:"Header"`
	Footer  Instructions `xml:"Footer"`
	Body    Instructions `xml:"Body"`
}

type Meta struct {
	XMLName xml.Name `xml:"Meta"`
	Author  string   `xml:"Author"`
	Creator string   `xml:"Creator"`
	Subject string   `xml:"Subject"`
}

type Default struct {
	XMLName     xml.Name      `xml:"Default"`
	Orientation Orientation   `xml:"Orientation"`
	Unit        Unit          `xml:"Unit"`
	Format      Format        `xml:"Format"`
	PageBreaks  PageBreakMode `xml:"PageBreaks"`
	PageMargins PageMargins   `xml:"PageMargins"`
}

func (scs *StyleClasses) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.CharData:
			pscs, err := ParseClasses(tok)
			if err != nil {
				return err
			}
			*scs = pscs
		case xml.EndElement:
			if tok == start.End() {
				return nil
			}
		default:
			return errors.Errorf("invalid xml token type for style element (%T)", tok)
		}
	}
}

func (is *Instructions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
			*is = append(*is, i)
		}
	}
}
