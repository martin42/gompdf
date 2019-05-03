package gompdf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/martin42/gompdf/style"

	"github.com/pkg/errors"
)

func ParseAndBuild(source string, target string) error {
	start := time.Now()
	fmt.Printf("load (%s) ...\n", source)
	doc, err := LoadFromFile(source)
	if err != nil {
		return err
	}
	fmt.Printf("loaded (%s) ... in (%s)\n", source, time.Since(start))

	outF, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outF.Close()

	start = time.Now()
	fmt.Printf("process ...\n")
	p, err := NewProcessor(doc)
	if err != nil {
		return err
	}
	err = p.Process(outF)
	if err != nil {
		return err
	}
	fmt.Printf("processed ... in (%s)\n", time.Since(start))
	return nil
}

func Load(r io.Reader) (*Document, error) {
	doc := &Document{}
	err := xml.NewDecoder(r).Decode(doc)
	if err != nil {
		return nil, err
	}
	doc.styleClasses, err = style.DecodeClasses(bytes.NewBufferString(doc.Style))
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

func (doc *Document) StyleClasses() style.Classes {
	return doc.styleClasses
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
	XMLName      xml.Name `xml:"Document"`
	Meta         Meta     `xml:"Meta"`
	Default      Default  `xml:"Default"`
	Style        string   `xml:"Style"`
	styleClasses style.Classes
	Header       Instructions `xml:"Header"`
	Footer       Instructions `xml:"Footer"`
	Body         Instructions `xml:"Body"`
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

func (is *Instructions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	err := is.DecodeAttrs(start.Attr)
	if err != nil {
		return err
	}
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
			is.iss = append(is.iss, i)
		}
	}
}
