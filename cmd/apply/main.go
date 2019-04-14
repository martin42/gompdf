package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/martin42/gompdf/style"
)

func main() {
	st := `
	buddy {
		width: 32.14;
		height: 14.56;
		color: #ff0088;
	}

	holly {
		font-point-size: 14.5;
		h-align: center;
		color: #112233;
	}
	`

	styles := style.Styles{}
	cs, err := style.DecodeClasses(bytes.NewBufferString(st))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	cs.Apply(&styles, "buddy", "holly")
	b, _ := json.MarshalIndent(styles, "", "  ")
	fmt.Printf("styles:\n%s\n", b)

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

	color: #ff0088;
	background-color: #995522;
	`

	// s := `
	// font-family: arial;
	// border: 1, 2, x, 4;
	// `

	//styles := style.Styles{}

	start := time.Now()
	app, err := style.DecodeApplier(bytes.NewBufferString(s))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("decoded in: %s\n", time.Since(start))
		start = time.Now()
		app.Apply(&styles)
		fmt.Printf("applied in: %s\n", time.Since(start))
		b, _ := json.MarshalIndent(styles, "", "  ")
		fmt.Printf("styles:\n%s\n", b)
	}
}
