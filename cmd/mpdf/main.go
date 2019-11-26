package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mazzegi/gompdf"
)

func main() {
	source := flag.String("source", "../../samples/doc1.xml", "")
	target := flag.String("target", "doc1.pdf", "")
	flag.Parse()

	fmt.Printf("compile (%s) to (%s) ...\n", *source, *target)
	start := time.Now()
	err := gompdf.ParseAndBuild(*source, *target)
	if err != nil {
		fmt.Printf("compile (%s) to (%s) ...failed: %v\n", *source, *target, err)
	} else {
		fmt.Printf("compile (%s) to (%s) ... done in (%s)\n", *source, *target, time.Since(start))
	}
}
