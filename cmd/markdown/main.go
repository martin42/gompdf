package main

import (
	"fmt"

	"github.com/mazzegi/gompdf/markdown"
)

func main() {
	tr := func(s string) string { return s }
	text := "some tricky ways to cut-off your tail"
	mdItems := markdown.NewProcessor().Process(text).WordItems(tr)
	dumpItems(text, mdItems)

	text = "*some* tricky `__ways__` to cut-off your tail\\"
	mdItems = markdown.NewProcessor().Process(text).WordItems(tr)
	dumpItems(text, mdItems)
}

func dumpItems(src string, items markdown.Items) {
	fmt.Printf("\nsrc: %q\n", src)
	for i, item := range items {
		fmt.Printf("%02d: %s\n", i, item)
	}
}
