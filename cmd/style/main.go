package main

import (
	"encoding/json"
	"fmt"

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
	`
	styles, err := style.NewDecoder().Decode(s)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		b, _ := json.MarshalIndent(styles, "", "  ")
		fmt.Printf("styles:\n%s\n", b)
	}
}
