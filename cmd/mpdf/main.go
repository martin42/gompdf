package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mazzegi/gompdf"
)

func main() {
	source := flag.String("source", "../../samples/doc6.xml", "")
	target := flag.String("target", "doc6.pdf", "")
	flag.Parse()
	build(*source, *target)
	//build(*source, "doc_copy.pdf")
}

func build(source, target string) {
	fmt.Printf("compile (%s) to (%s) ...\n", source, target)
	start := time.Now()
	err := gompdf.ParseAndBuildFile(source, target)
	if err != nil {
		fmt.Printf("compile (%s) to (%s) ...failed: %v\n", source, target, err)
	} else {
		fmt.Printf("compile (%s) to (%s) ... done in (%s)\n", source, target, time.Since(start))
	}
}
