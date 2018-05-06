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

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("unable to initialize sdl: %v", err)
	}
	defer sdl.Quit()
	sdl.GLSetSwapInterval(1)

	mach := mach85.New()
	if err := mach.Init(); err != nil {
		log.Fatalf("unable to initialize: %v", err)
	}
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
