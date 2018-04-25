package main

import (
	"encoding/json"
	"log"
	"os"

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

	in, err := os.Open("c64rom.debug")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(in)
	comments := []mach85.Comment{}
	decoder.Decode(&comments)
	mon.LoadComments(comments)

	mon.Run()
}
