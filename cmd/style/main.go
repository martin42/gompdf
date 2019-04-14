package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/martin42/gompdf/style"
)

func main() {
	s := `
	font-family: arial;
	font-point-size: 12;
	font-style: italic;
	font-weight: bold;
	font-decoration: underline;

	border: 1, 2, 3, 4;
	padding: 5.1, 5.2, 5.3, 5.4;
	margin: 7.1, 7.2, 7.3, 7.4;

	width: 64.8;
	height: 98.1;

	h-align: left;
	v-align: top;
	`

	start := time.Now()
	styles := style.Styles{
		Align: style.Align{
			HAlign: style.HAlignCenter,
			VAlign: style.VAlignBottom,
		},
	}

	err := style.NewDecoder(bytes.NewBufferString(s)).Decode(&styles)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("decoded in: %s\n", time.Since(start))
		b, _ := json.MarshalIndent(styles, "", "  ")
		fmt.Printf("styles:\n%s\n", b)
	}
}
