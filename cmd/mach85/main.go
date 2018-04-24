package main

import (
	"log"

	"github.com/blackchip-org/mach85"
)

func main() {
	log.SetFlags(0)

	log.Printf("\nWelcome to the Mach-85!\n\n")

	mach := mach85.New()
	if err := mach.LoadROM(); err != nil {
		log.Fatal(err)
	}
	log.Println()
	mon := mach85.NewMonitor(mach)
	mon.Run()
}
