package main

import (
	"flag"
	"fmt"

	"github.com/blackchip-org/mach85/encoding"
	"github.com/blackchip-org/mach85/encoding/petscii"
)

var shifted bool

func init() {
	flag.BoolVar(&shifted, "s", false, "show shifted PETSCII")
}

func main() {
	flag.Parse()
	var d encoding.Decoder
	if shifted {
		d = petscii.ShiftedDecoder{}
	} else {
		d = petscii.UnshiftedDecoder{}
	}
	fmt.Print("    ")
	for x := 0; x < 0x10; x++ {
		fmt.Printf("%x ", x)
	}
	fmt.Println()
	for y := 0; y < 0x10; y++ {
		fmt.Printf("%x   ", y)
		for x := 0; x < 0x10; x++ {
			char := uint8(y<<4 + x)
			if !d.IsPrintable(char) {
				fmt.Printf("  ")
			} else {
				fmt.Printf("%c ", d.Decode(char))
			}
		}
		fmt.Println()
	}
}
