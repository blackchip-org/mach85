package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/blackchip-org/mach85"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	log.SetFlags(0)

	log.Printf("\nWelcome to the Mach-85!\n\n")

	mach := mach85.New()
	if err := mach.LoadROM(); err != nil {
		log.Fatal(err)
	}
	log.Println()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("unable to initialize sdl: %v", err)
	}
	defer sdl.Quit()
	sdl.GLSetSwapInterval(1)
	video, err := mach85.NewVideo(mach.Memory)
	if err != nil {
		log.Fatalf("unable to create window: %v", err)
	}

	mach.AddDevice(video)
	mach.AddDevice(mach85.NewHackDevice(mach.Memory))

	mon := mach85.NewMonitor(mach)
	in, err := os.Open("c64rom_en.source")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(in)
	source := &mach85.Source{}
	decoder.Decode(source)
	mon.Disassembler.LoadSource(source)

	mon.Run()
}
