package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/martin42/gompdf"
)

func main() {
	source := flag.String("source", "../../samples/doc2.xml", "")
	target := flag.String("target", "doc2.pdf", "")
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
