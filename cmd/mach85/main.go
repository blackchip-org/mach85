package main

import (
	"log"

	"github.com/blackchip-org/mach85"
)

func main() {
	mach := mach85.New()
	if err := mach.LoadROM(); err != nil {
		log.Fatal(err)
	}
	mon := mach85.NewMonitor(mach)
	mon.Run()
}
