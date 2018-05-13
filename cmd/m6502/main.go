package main

import (
	"log"

	"github.com/blackchip-org/mach85"
)

func main() {
	log.SetFlags(0)

	mach := mach85.New()
	mon := mach85.NewMonitor(mach)
	mon.Prompt = "m6502"
	mon.Run()
}
