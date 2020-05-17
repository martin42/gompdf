package gompdf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/mazzegi/gompdf/style"

	"github.com/pkg/errors"
)

func ParseAndBuild(source io.Reader, target io.Writer) error {
	doc, err := Load(source)
	if err != nil {
		return err
	}
	start := time.Now()
	fmt.Printf("process first time...\n")
	p, err := NewProcessor(doc)
	if err != nil {
		return err
	}
	err = p.Process(ioutil.Discard, 0)
	if err != nil {
		return err
	}
	numPages := p.currPage
	fmt.Printf("processed first time... in (%s)\n", time.Since(start))

	start = time.Now()
	fmt.Printf("process second time...\n")
	p, err = NewProcessor(doc)
	if err != nil {
		return err
	}
	err = p.Process(target, numPages)
	if err != nil {
		return err
	}
	fmt.Printf("processed second time... in (%s)\n", time.Since(start))

	return nil
}

func ParseAndBuildFile(source string, target string) error {
	srcF, err := os.Open(source)
	if err != nil {
		return errors.Errorf("open (%s)", source)
	}
	defer srcF.Close()

	outF, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outF.Close()

	return ParseAndBuild(srcF, outF)
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
	XMLName xml.Name `xml:"page-margins"`
	Left    float64  `xml:"left"`
	Top     float64  `xml:"top"`
	Right   float64  `xml:"right"`
	Bottom  float64  `xml:"bottom"`
}

type Document struct {
	XMLName      xml.Name `xml:"document"`
	Meta         Meta     `xml:"meta"`
	Default      Default  `xml:"default"`
	Style        string   `xml:"style"`
	styleClasses style.Classes
	Header       Instructions `xml:"header"`
	Footer       Instructions `xml:"footer"`
	Body         Instructions `xml:"body"`
}

type Meta struct {
	XMLName xml.Name `xml:"meta"`
	Author  string   `xml:"author"`
	Creator string   `xml:"creator"`
	Subject string   `xml:"subject"`
}

type Default struct {
	XMLName     xml.Name      `xml:"default"`
	Orientation Orientation   `xml:"orientation"`
	Unit        Unit          `xml:"unit"`
	Format      Format        `xml:"format"`
	PageBreaks  PageBreakMode `xml:"page-breaks"`
	PageMargins PageMargins   `xml:"page-margins"`
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
