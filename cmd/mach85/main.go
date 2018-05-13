package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blackchip-org/mach85"
	"github.com/blackchip-org/mach85/rom"
	"github.com/veandco/go-sdl2/sdl"
)

var wait bool

func init() {
	flag.BoolVar(&wait, "w", false, "wait for user to issue go command")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

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
	in, err := os.Open(filepath.Join(rom.Path, "c64rom_en.source"))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(in)
	source := &mach85.Source{}
	decoder.Decode(source)
	mon.Disassembler.LoadSource(source)

	if wait {
		mon.Run()
	} else {
		fmt.Println("Press control-C to start monitor")
		mon.Go()
	}
}
