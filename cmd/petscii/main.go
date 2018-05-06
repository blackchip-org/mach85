package main

import (
	"fmt"

	"github.com/blackchip-org/mach85"
)

func main() {
	fmt.Printf("unshifted\n")
	printTable(mach85.PetsciiUnshiftedDecoder)
	fmt.Printf("\nshifted\n")
	printTable(mach85.PetsciiShiftedDecoder)
	fmt.Printf("\nscreen unshifted\n")
	printTable(mach85.ScreenUnshiftedDecoder)
	fmt.Printf("\nscreen shifted\n")
	printTable(mach85.ScreenShiftedDecoder)
}

func printTable(decode mach85.Decoder) {
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
