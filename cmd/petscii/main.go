package main

import (
	"fmt"

	"github.com/blackchip-org/mach85/encoding"
	"github.com/blackchip-org/mach85/encoding/petscii"
)

func main() {
	fmt.Printf("unshifted\n")
	printTable(petscii.UnshiftedDecoder)
	fmt.Printf("\nshifted\n")
	printTable(petscii.ShiftedDecoder)
}

func printTable(decode encoding.Decoder) {
	fmt.Print("    ")
	for x := 0; x < 0x10; x++ {
		fmt.Printf("%x ", x)
	}
	fmt.Println()
	for y := 0; y < 0x10; y++ {
		fmt.Printf("%x   ", y)
		for x := 0; x < 0x10; x++ {
			code := uint8(y<<4 + x)
			ch, printable := decode(code)
			if !printable {
				fmt.Printf("  ")
			} else {
				fmt.Printf("%c ", ch)
			}
		}
		fmt.Println()
	}
}
