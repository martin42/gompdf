package main

import "fmt"

func main() {
	fmt.Printf("Hallo ich bin Leni, ich bin Gott ;-)\n")

	for {
		fmt.Printf("bitte eine Zahl eingeben: ")
		var z1 int
		fmt.Scanln(&z1)

		fmt.Printf("bitte eine weitere Zahl eingeben: ")
		var z2 int
		fmt.Scanln(&z2)

		fmt.Printf("%d + %d = %d\n\n", z1, z2, z1+z2)
	}

}
