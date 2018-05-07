package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

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

	file := filepath.Join("..", "..", "rom", "chargen.rom")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	chargen := mach85.NewROM(data)
	w := 32 * 8
	h := 16 * 8
	scale := 4

	window, err := sdl.CreateWindow(
		"chargen",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(w*scale), int32(h*scale),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("unable to initialize window: %v", err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("unable to initialize renderer: %v", err)
	}

	clear := sdl.Rect{X: 0, Y: 0, W: int32(w * scale), H: int32(h * scale)}
	renderer.SetScale(float32(scale), float32(scale))
	renderer.SetDrawColorArray(mach85.Blue...)
	renderer.FillRect(&clear)

	renderer.SetDrawColorArray(mach85.LightBlue...)
	baseX := 0
	baseY := 0
	addr := uint16(0)
	for baseY < h {
		for y := baseY; y < baseY+8; y++ {
			line := chargen.Load(addr)
			addr++
			for x := baseX; x < baseX+8; x++ {
				bit := line & 0x80
				line = line << 1
				if bit != 0 {
					renderer.DrawPoint(int32(x), int32(y))
				}
			}
		}
		baseX += 8
		if baseX >= w {
			baseX = 0
			baseY += 8
		}
	}

	renderer.Present()

	run := true
	for run {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				run = false
			}
		}
	}
}
