package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mazzegi/gompdf"
	"github.com/signintech/gopdf"
)

func main() {
	source := flag.String("source", "../../samples/simple.xml", "")
	target := flag.String("target", "simple.pdf", "")
	flag.Parse()
	build(*source, *target)
	//build(*source, "doc_copy.pdf")
}

func build(source, target string) {
	e := goPdfEngine()
	fmt.Printf("compile (%s) to (%s) ...\n", source, target)
	start := time.Now()
	err := gompdf.ParseAndBuildFile(e, source, target)
	if err != nil {
		fmt.Printf("compile (%s) to (%s) ...failed: %v\n", source, target, err)
	} else {
		fmt.Printf("compile (%s) to (%s) ... done in (%s)\n", source, target, time.Since(start))
	}
}

func goPdfEngine() *gompdf.GoPdfEngine {

	regularOption := gopdf.TtfOption{
		Style: gopdf.Regular,
	}
	boldOption := gopdf.TtfOption{
		Style: gopdf.Bold,
	}
	italicOption := gopdf.TtfOption{
		Style: gopdf.Italic,
	}
	boldItalicOption := gopdf.TtfOption{
		Style: gopdf.Bold | gopdf.Italic,
	}

	fonts := []gompdf.GoPdfTTFFont{}
	add := func(name string, ttfPath string, option gopdf.TtfOption) {
		fonts = append(fonts, gompdf.GoPdfTTFFont{
			Name:    name,
			TTFPath: ttfPath,
			Option:  option,
		})
	}

	add("sans", "fonts/NotoSans-Regular.ttf", regularOption)
	add("sans", "fonts/NotoSans-Bold.ttf", boldOption)
	add("sans", "fonts/NotoSans-Italic.ttf", italicOption)
	add("sans", "fonts/NotoSans-BoldItalic.ttf", boldItalicOption)
	add("mono", "fonts/LiberationMono-Regular.ttf", regularOption)
	add("mono", "fonts/LiberationMono-Bold.ttf", boldOption)
	add("mono", "fonts/LiberationMono-Italic.ttf", italicOption)
	add("mono", "fonts/LiberationMono-BoldItalic.ttf", boldItalicOption)
	p := gompdf.NewGoPdfEngine(fonts...)
	return p
}
