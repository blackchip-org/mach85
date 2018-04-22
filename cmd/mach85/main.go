package main

import (
	"log"

	"github.com/blackchip-org/mach85"
)

func main() {
	mach, err := mach85.New()
	if err != nil {
		log.Fatal(err)
	}
	mach.Run()
}
